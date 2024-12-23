package usecase

import (
	"context"
	"fmt"

	"github.com/idoyudha/eshop-warehouse/internal/entity"
	"github.com/idoyudha/eshop-warehouse/internal/utils"
)

type TransactionProductUseCase struct {
	repoRedis              WarehouseRankRedisRepo
	repoTransactionPostgre TransactionProductPostgresRepo
	repoProductPostgre     WarehouseProductPostgreRepo
}

func NewTransactionProductUseCase(
	repoRedis WarehouseRankRedisRepo,
	repoTransactionPostgre TransactionProductPostgresRepo,
	repoProductPostgre WarehouseProductPostgreRepo,
) *TransactionProductUseCase {
	return &TransactionProductUseCase{
		repoRedis,
		repoTransactionPostgre,
		repoProductPostgre,
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

// move from warehouse to user
func (u *TransactionProductUseCase) MoveOut(ctx context.Context, stockMovementReq []*entity.StockMovement, zipCode string) error {
	var stockMovements []*entity.StockMovement
	for _, stockMovement := range stockMovementReq {
		// find the nearest warehouse and remaining product quantity for each warehouse
		// TODO: can be improved if we get from in memory database (redis)
		warehouses, err := u.repoProductPostgre.GetWarehouseIDZipCodeAndQtyByProductID(ctx, stockMovement.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get warehouse id and zip code by product id: %w", err)
		}

		var totalProductQuantity int64
		for _, warehouse := range warehouses {
			totalProductQuantity += warehouse.ProductQuantity
		}
		if totalProductQuantity < stockMovement.Quantity {
			return fmt.Errorf("%s requested product quantity is not enough", warehouses[0].ProductName)
		}

		nearestWarehouseIDs, err := utils.FindNearestWarehouse(zipCode, warehouses, stockMovement.Quantity)
		if err != nil {
			return fmt.Errorf("failed to calculate nearest warehouse: %w", err)
		}

		for warehouseID, quantity := range nearestWarehouseIDs {
			stockMovements = append(stockMovements, &entity.StockMovement{
				ID:              stockMovement.ID,
				ProductID:       stockMovement.ProductID,
				ProductName:     stockMovement.ProductName,
				Quantity:        quantity,
				FromWarehouseID: warehouseID,
				ToUserID:        stockMovement.ToUserID,
				CreatedAt:       stockMovement.CreatedAt,
			})
		}
	}

	return u.repoTransactionPostgre.TransferOut(ctx, stockMovements)
}
