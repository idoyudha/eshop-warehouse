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
		GetAllExceptMain(ctx context.Context) ([]*entity.Warehouse, error)
	}

	WarehouseProductPostgreRepo interface {
		Save(ctx context.Context, warehouseProduct *entity.WarehouseProduct) error
		UpdateProductQuantity(ctx context.Context, warehouseProduct *entity.WarehouseProduct) error
		GetAll(ctx context.Context) ([]*entity.WarehouseProduct, error)
		GetByProductID(ctx context.Context, id uuid.UUID) (*entity.WarehouseProduct, error)
		GetByWarehouseID(ctx context.Context, id uuid.UUID) ([]*entity.WarehouseProduct, error)
		GetByProductIDAndWarehouseID(ctx context.Context, productID uuid.UUID, warehouseID uuid.UUID) (*entity.WarehouseProduct, error)
	}

	StockMovementPostgreRepo interface {
		GetAll(ctx context.Context) ([]*entity.StockMovement, error)
		GetByProductID(context.Context, uuid.UUID) ([]*entity.StockMovement, error)
		GetBySourceID(context.Context, uuid.UUID) ([]*entity.StockMovement, error)
		GetByDestinationID(context.Context, uuid.UUID) ([]*entity.StockMovement, error)
	}

	TransactionProductPostgresRepo interface {
		TransferIn(ctx context.Context, stockMovement *entity.StockMovement, warehouseProduct *entity.WarehouseProduct) error
		TransferOut(ctx context.Context, stockMovement []*entity.StockMovement, warehouseProduct *entity.WarehouseProduct) error
	}

	Warehouse interface {
		CreateWarehouse(ctx context.Context, warehouse *entity.Warehouse) error
		UpdateWarehouse(ctx context.Context, warehouse *entity.Warehouse) error
		GetWarehouseByID(ctx context.Context, id uuid.UUID) (*entity.Warehouse, error)
		GetAllWarehouses(ctx context.Context) ([]*entity.Warehouse, error)
	}

	WarehouseProduct interface {
		CreateWarehouseProduct(ctx context.Context, warehouseProduct *entity.WarehouseProduct) error
		UpdateWarehouseProduct(ctx context.Context, warehouseProduct *entity.WarehouseProduct) error
		GetAllWarehouseProducts(ctx context.Context) ([]*entity.WarehouseProduct, error)
		GetWarehouseProductByProductID(ctx context.Context, id uuid.UUID) (*entity.WarehouseProduct, error)
		GetWarehouseProductByWarehouseID(ctx context.Context, id uuid.UUID) ([]*entity.WarehouseProduct, error)
		GetWarehouseProductByProductIDAndWarehouseID(ctx context.Context, productID uuid.UUID, warehouseID uuid.UUID) (*entity.WarehouseProduct, error)
	}

	StockMovement interface {
		MoveIn(context.Context, uuid.UUID, uuid.UUID, *entity.WarehouseProduct) error
		MoveOut(context.Context, uuid.UUID, uuid.UUID, *entity.WarehouseProduct) error
		GetMovemenetByProductID(context.Context, uuid.UUID) ([]*entity.StockMovement, error)
		GetMovementBySourceID(context.Context, uuid.UUID) ([]*entity.StockMovement, error)
		GetMovementByDestinationID(context.Context, uuid.UUID) ([]*entity.StockMovement, error)
	}

	TransactionProductMovement interface {
		MoveIn(context.Context, *entity.StockMovement) error
		MoveOut(context.Context, []*entity.StockMovement) error
	}
)
