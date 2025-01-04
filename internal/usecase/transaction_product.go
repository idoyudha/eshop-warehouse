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

		// TODO: update product stock in product service
	}

	return u.repoTransactionPostgre.TransferOut(ctx, stockMovements)
}
