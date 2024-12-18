package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/entity"
)

type (
	WarehousePostgreRepo interface {
		Save(ctx context.Context, warehouse *entity.Warehouse) error
		Update(ctx context.Context, warehouse *entity.Warehouse) error
		GetByID(ctx context.Context, id uuid.UUID) (*entity.Warehouse, error)
		GetAll(ctx context.Context) ([]*entity.Warehouse, error)
	}

	WarehouseProductPostgreRepo interface {
		Save(ctx context.Context, warehouseProduct *entity.WarehouseProduct) error
		Update(ctx context.Context, warehouseProduct *entity.WarehouseProduct) error
		GetByProductID(ctx context.Context, id uuid.UUID) (*entity.WarehouseProduct, error)
		GetByWarehouseID(ctx context.Context, id uuid.UUID) ([]*entity.WarehouseProduct, error)
		GetByProductIDAndWarehouseID(ctx context.Context, productID uuid.UUID, warehouseID uuid.UUID) (*entity.WarehouseProduct, error)
	}

	StockMovementPostgreRepo interface {
		TransferIn(context.Context, uuid.UUID, uuid.UUID, *entity.WarehouseProduct) error
		TrasnferOut(context.Context, uuid.UUID, uuid.UUID, *entity.WarehouseProduct) error
		GetByProductID(context.Context, uuid.UUID) ([]*entity.StockMovement, error)
		GetBySourceID(context.Context, uuid.UUID) ([]*entity.StockMovement, error)
		GetByDestinationID(context.Context, uuid.UUID) ([]*entity.StockMovement, error)
	}

	Warehouse interface {
		CreateWarehouse(ctx context.Context, warehouse *entity.Warehouse) error
		UpdateWarehouse(ctx context.Context, warehouse *entity.Warehouse) error
		GetWarehouseByID(ctx context.Context, id uuid.UUID) (*entity.Warehouse, error)
		GetAllWarehouses(ctx context.Context) ([]*entity.Warehouse, error)
	}

	WarehouseProductInterface interface {
		CreateWarehouseProduct(ctx context.Context, warehouseProduct *entity.WarehouseProduct) error
		UpdateWarehouseProduct(ctx context.Context, warehouseProduct *entity.WarehouseProduct) error
		GetAllWarehouseProducts(ctx context.Context) ([]*entity.WarehouseProduct, error)
		GetWarehouseProductByProductID(ctx context.Context, id uuid.UUID) (*entity.WarehouseProduct, error)
		GetWarehouseProductByWarehouseID(ctx context.Context, id uuid.UUID) ([]*entity.WarehouseProduct, error)
		GetWarehouseProductByProductIDAndWarehouseID(ctx context.Context, productID uuid.UUID, warehouseID uuid.UUID) (*entity.WarehouseProduct, error)
	}

	StockMovementInterface interface {
		MoveIn(context.Context, uuid.UUID, uuid.UUID, *entity.WarehouseProduct) error
		MoveOut(context.Context, uuid.UUID, uuid.UUID, *entity.WarehouseProduct) error
		GetMovemenetByProductID(context.Context, uuid.UUID) ([]*entity.StockMovement, error)
		GetMovementBySourceID(context.Context, uuid.UUID) ([]*entity.StockMovement, error)
		GetMovementByDestinationID(context.Context, uuid.UUID) ([]*entity.StockMovement, error)
	}
)
