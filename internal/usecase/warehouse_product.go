package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/entity"
)

type WarehouseProductUseCase struct {
	repoPostgre WarehouseProductPostgreRepo
}

func NewWarehouseProductUseCase(repoPostgre WarehouseProductPostgreRepo) *WarehouseProductUseCase {
	return &WarehouseProductUseCase{
		repoPostgre,
	}
}

func (u *WarehouseProductUseCase) CreateWarehouseProduct(ctx context.Context, warehouseProduct *entity.WarehouseProduct) error {
	return u.repoPostgre.Save(ctx, warehouseProduct)
}

func (u *WarehouseProductUseCase) UpdateWarehouseProductNameAndPrice(ctx context.Context, warehouseProduct *entity.WarehouseProduct) error {
	return u.repoPostgre.UpdateNameAndPrice(ctx, warehouseProduct)
}

// TODO: need to handle stock movement also if quantity is updated
func (u *WarehouseProductUseCase) UpdateWarehouseProductQuantity(ctx context.Context, warehouseProduct *entity.WarehouseProduct) error {
	return u.repoPostgre.UpdateProductQuantity(ctx, warehouseProduct)
}

func (u *WarehouseProductUseCase) GetAllWarehouseProducts(ctx context.Context) ([]*entity.WarehouseProduct, error) {
	return u.repoPostgre.GetAll(ctx)
}

func (u *WarehouseProductUseCase) GetWarehouseProductByProductID(ctx context.Context, productID uuid.UUID) ([]*entity.WarehouseProduct, error) {
	return u.repoPostgre.GetByProductID(ctx, productID)
}

func (u *WarehouseProductUseCase) GetWarehouseProductByWarehouseID(ctx context.Context, warehouseID uuid.UUID) ([]*entity.WarehouseProduct, error) {
	return u.repoPostgre.GetByWarehouseID(ctx, warehouseID)
}

func (u *WarehouseProductUseCase) GetWarehouseProductByProductIDAndWarehouseID(ctx context.Context, productID uuid.UUID, warehouseID uuid.UUID) (*entity.WarehouseProduct, error) {
	return u.repoPostgre.GetByProductIDAndWarehouseID(ctx, productID, warehouseID)
}
