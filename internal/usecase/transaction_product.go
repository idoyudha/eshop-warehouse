package usecase

import (
	"context"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/entity"
)

type TransactionProductUseCase struct {
	repoRedis              WarehouseRankRedisRepo
	repoTransactionPostgre TransactionProductPostgresRepo
	repoProductPostgre     WarehouseProductPostgreRepo
}

func NewTransactionProductUseCase(repoRedis WarehouseRankRedisRepo, repoTransactionPostgre TransactionProductPostgresRepo, repoProductPostgre WarehouseProductPostgreRepo) *TransactionProductUseCase {
	return &TransactionProductUseCase{
		repoRedis,
		repoTransactionPostgre,
		repoProductPostgre,
	}
}

func (u *TransactionProductUseCase) MoveIn(ctx context.Context, stockMovement *entity.StockMovement) error {
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
func (u *TransactionProductUseCase) MoveOut(ctx context.Context, stockMovement *entity.StockMovement) error {
	// check if product quantity is enough for all warehouses
	products, err := u.repoProductPostgre.GetByProductID(ctx, stockMovement.ProductID)
	if err != nil {
		return err
	}
	var totalProductQuantity int64
	for _, product := range products {
		totalProductQuantity += product.ProductQuantity
	}
	if totalProductQuantity < stockMovement.Quantity {
		return fmt.Errorf("product quantity is not enough")
	}

	// get product stock
	product, err := u.repoProductPostgre.GetByProductIDAndWarehouseID(ctx, stockMovement.ProductID, stockMovement.FromWarehouseID)
	if err != nil {
		return err
	}

	var stockMovements []*entity.StockMovement
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	quantityFromInitial := math.Min(float64(product.ProductQuantity), float64(product.ProductQuantity))
	curStockMovement := &entity.StockMovement{
		ID:              id,
		ProductID:       stockMovement.ProductID,
		ProductName:     stockMovement.ProductName,
		Quantity:        product.ProductQuantity,
		FromWarehouseID: stockMovement.FromWarehouseID,
		ToUserID:        stockMovement.ToUserID,
		CreatedAt:       stockMovement.CreatedAt,
	}
	stockMovements = append(stockMovements, curStockMovement)
	remainingQuantity := stockMovement.Quantity - int64(quantityFromInitial)

	if remainingQuantity > 0 {
		// 1. get from nearest warehouse -> using redis rank
		nearestWarehouses, err := u.repoRedis.GetNearestWarehouses(ctx, stockMovement.FromWarehouseID.String())
		if err != nil {
			return err
		}

		// 2. get product quantity from nearest warehouse
		for _, nearWarehouse := range nearestWarehouses {
			if remainingQuantity <= 0 {
				break
			}

			// skip inital warehouse
			if nearWarehouse.WarehouseID == stockMovement.FromWarehouseID.String() {
				continue
			}

			warehouseUUID := uuid.MustParse(nearWarehouse.WarehouseID)
			nearestProduct, err := u.repoProductPostgre.GetByProductIDAndWarehouseID(ctx, stockMovement.ProductID, warehouseUUID)
			if err != nil {
				return err
			}

			// if nearest is empty, skip
			if nearestProduct.ProductQuantity <= 0 {
				continue
			}

			quantityFromWarehouse := math.Min(float64(nearestProduct.ProductQuantity), float64(remainingQuantity))

			id, err := uuid.NewV7()
			if err != nil {
				return err
			}
			tmpStockMovement := &entity.StockMovement{
				ID:              id,
				ProductID:       stockMovement.ProductID,
				ProductName:     stockMovement.ProductName,
				Quantity:        int64(quantityFromWarehouse),
				FromWarehouseID: warehouseUUID,
				ToUserID:        stockMovement.ToUserID,
				CreatedAt:       stockMovement.CreatedAt,
			}
			stockMovements = append(stockMovements, tmpStockMovement)
			remainingQuantity -= int64(quantityFromWarehouse)
		}
	}
	return u.repoTransactionPostgre.TransferOut(ctx, stockMovements)
}
