package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/entity"
	"github.com/idoyudha/eshop-warehouse/internal/utils"
	"github.com/idoyudha/eshop-warehouse/pkg/kafka"
)

const productQuantityUpdated = "product-quantity-updated"

type TransactionProductUseCase struct {
	repoRedis              WarehouseRankRedisRepo
	repoTransactionPostgre TransactionProductPostgresRepo
	repoProductPostgre     WarehouseProductPostgreRepo
	producer               *kafka.ProducerServer
}

func NewTransactionProductUseCase(
	repoRedis WarehouseRankRedisRepo,
	repoTransactionPostgre TransactionProductPostgresRepo,
	repoProductPostgre WarehouseProductPostgreRepo,
	producer *kafka.ProducerServer,
) *TransactionProductUseCase {
	return &TransactionProductUseCase{
		repoRedis,
		repoTransactionPostgre,
		repoProductPostgre,
		producer,
	}
}

func (u *TransactionProductUseCase) MoveIn(ctx context.Context, stockMovement *entity.StockMovement) error {
	err := stockMovement.GenerateStockMovementID()
	if err != nil {
		return err
	}

	warehouseProduct, err := u.repoProductPostgre.GetByProductIDAndWarehouseID(ctx, stockMovement.ProductID, stockMovement.FromWarehouseID)
	if err != nil {
		return err
	}

	if warehouseProduct.ProductQuantity < stockMovement.Quantity {
		return fmt.Errorf("warehouse product quantity is not enough")
	}

	return u.repoTransactionPostgre.TransferIn(ctx, stockMovement)
}

type kafkaProductQuantityUpdatedMessage struct {
	ProductID uuid.UUID
	Quantity  int
}

// move from warehouse to user
func (u *TransactionProductUseCase) MoveOut(ctx context.Context, stockMovementReq []*entity.StockMovement, zipCode string) error {
	var stockMovements []*entity.StockMovement
	for _, stockMovement := range stockMovementReq {
		totalProduct, err := u.repoProductPostgre.GetTotalQuantityOfProductInAllWarehouse(ctx, stockMovement.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get total quantity of product in all warehouse: %w", err)
		}

		if totalProduct < int(stockMovement.Quantity) {
			return fmt.Errorf("product quantity is not enough")
		}
		// find the nearest warehouse and remaining product quantity for each warehouse
		// TODO: can be improved if we get from in memory database (redis)
		warehouses, err := u.repoProductPostgre.GetWarehouseIDZipCodeAndQtyByProductID(ctx, stockMovement.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get warehouse id and zip code by product id: %w", err)
		}
		nearestWarehouseIDs, err := utils.FindNearestWarehouse(zipCode, warehouses, stockMovement.Quantity)
		if err != nil {
			return fmt.Errorf("failed to calculate nearest warehouse: %w", err)
		}

		for warehouseID, quantity := range nearestWarehouseIDs {
			var newStockMovement entity.StockMovement
			err = newStockMovement.GenerateStockMovementID()
			if err != nil {
				return fmt.Errorf("failed to generate stock movement id: %w", err)
			}
			newStockMovement.ProductID = stockMovement.ProductID
			newStockMovement.ProductName = stockMovement.ProductName
			newStockMovement.Quantity = quantity
			newStockMovement.FromWarehouseID = warehouseID
			newStockMovement.ToUserID = stockMovement.ToUserID
			newStockMovement.CreatedAt = stockMovement.CreatedAt
			stockMovements = append(stockMovements, &newStockMovement)
		}

		message := kafkaProductQuantityUpdatedMessage{
			ProductID: stockMovement.ProductID,
			Quantity:  totalProduct - int(stockMovement.Quantity),
		}

		err = u.producer.Produce(
			productQuantityUpdated,
			[]byte(stockMovement.ProductID.String()),
			message,
		)
		if err != nil {
			// TODO: handle error, cancel the update if failed. or try use retry mechanism
			return fmt.Errorf("failed to produce kafka message: %w", err)
		}
	}

	return u.repoTransactionPostgre.TransferOut(ctx, stockMovements)
}
