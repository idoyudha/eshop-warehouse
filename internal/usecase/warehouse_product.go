package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/entity"
	"github.com/idoyudha/eshop-warehouse/internal/utils"
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
	err := warehouseProduct.GenerateWarehouseProductID()
	if err != nil {
		return err
	}
	return u.repoPostgre.Save(ctx, warehouseProduct)
}

func (u *WarehouseProductUseCase) UpdateWarehouseProduct(ctx context.Context, warehouseProduct *entity.WarehouseProduct) error {
	return u.repoPostgre.Update(ctx, warehouseProduct)
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

func (u *WarehouseProductUseCase) GetNearestWarehouseZipCodeByProductID(ctx context.Context, zipCode string, productID uuid.UUID) (*string, error) {
	warehouse, err := u.repoPostgre.GetWarehouseIDZipCodeAndQtyByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get warehouse and zipcode data: %w", err)
	}

	zipCodeRes, err := utils.FindNearestWarehouseByProductID(zipCode, warehouse)
	if err != nil {
		return nil, fmt.Errorf("failed to find nearest warehouse: %w", err)
	}

	return &zipCodeRes, nil
}
