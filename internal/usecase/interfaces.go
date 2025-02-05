package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/entity"
)

type (
	WarehousePostgreRepo interface {
		Save(context.Context, *entity.Warehouse) error
		Update(context.Context, *entity.Warehouse) error
		GetByID(context.Context, uuid.UUID) (*entity.Warehouse, error)
		GetAll(context.Context) ([]*entity.Warehouse, error)
		GetAllExceptMain(context.Context) ([]*entity.Warehouse, error)
		GetMainID(context.Context) (uuid.UUID, error)
		GetAllIDAndZipCode(context.Context) ([]*entity.Warehouse, error)
	}

	WarehouseProductPostgreRepo interface {
		Save(context.Context, *entity.WarehouseProduct) error
		Update(context.Context, *entity.WarehouseProduct) error
		UpdateProductQuantity(context.Context, *entity.WarehouseProduct) error
		GetAll(context.Context) ([]*entity.WarehouseProduct, error)
		GetByProductID(context.Context, uuid.UUID) ([]*entity.WarehouseProduct, error)
		GetByWarehouseID(context.Context, uuid.UUID) ([]*entity.WarehouseProduct, error)
		GetByProductIDAndWarehouseID(context.Context, uuid.UUID, uuid.UUID) (*entity.WarehouseProduct, error)
		GetWarehouseIDZipCodeAndQtyByProductID(context.Context, uuid.UUID) ([]*entity.WarehouseAddressAndProductQty, error)
		GetTotalQuantityOfProductInAllWarehouse(context.Context, uuid.UUID) (int, error)
	}

	StockMovementPostgreRepo interface {
		GetAll(ctx context.Context) ([]*entity.StockMovement, error)
		GetByProductID(context.Context, uuid.UUID) ([]*entity.StockMovement, error)
		GetBySourceID(context.Context, uuid.UUID) ([]*entity.StockMovement, error)
		GetByDestinationID(context.Context, uuid.UUID) ([]*entity.StockMovement, error)
	}

	TransactionProductPostgresRepo interface {
		TransferIn(context.Context, *entity.StockMovement) error
		TransferOut(context.Context, []*entity.StockMovement) error
	}

	Warehouse interface {
		CreateWarehouse(context.Context, *entity.Warehouse) error
		UpdateWarehouse(context.Context, *entity.Warehouse) error
		GetWarehouseByID(context.Context, uuid.UUID) (*entity.Warehouse, error)
		GetAllWarehouses(context.Context) ([]*entity.Warehouse, error)
		GetMainIDWarehouse(context.Context) (uuid.UUID, error)
		GetNearestWarehouse(context.Context, []string) (map[string]string, error)
	}

	WarehouseProduct interface {
		CreateWarehouseProduct(context.Context, *entity.WarehouseProduct) error
		UpdateWarehouseProduct(context.Context, *entity.WarehouseProduct) error
		UpdateWarehouseProductQuantity(context.Context, *entity.WarehouseProduct) error
		GetAllWarehouseProducts(context.Context) ([]*entity.WarehouseProduct, error)
		GetWarehouseProductByProductID(context.Context, uuid.UUID) ([]*entity.WarehouseProduct, error)
		GetWarehouseProductByWarehouseID(context.Context, uuid.UUID) ([]*entity.WarehouseProduct, error)
		GetWarehouseProductByProductIDAndWarehouseID(context.Context, uuid.UUID, uuid.UUID) (*entity.WarehouseProduct, error)
		GetNearestWarehouseZipCodeByProductID(context.Context, string, uuid.UUID) (*string, error)
	}

	StockMovement interface {
		GetAllStockMovements(context.Context) ([]*entity.StockMovement, error)
		GetStockMovementsByProductID(context.Context, uuid.UUID) ([]*entity.StockMovement, error)
		GetStockMovementsBySourceID(context.Context, uuid.UUID) ([]*entity.StockMovement, error)
		GetStockMovementsByDestinationID(context.Context, uuid.UUID) ([]*entity.StockMovement, error)
	}

	TransactionProduct interface {
		MoveIn(context.Context, *entity.StockMovement) error
		MoveOut(context.Context, []*entity.StockMovement, string) error
	}
)
