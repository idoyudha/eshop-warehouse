package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/entity"
	"github.com/idoyudha/eshop-warehouse/internal/utils"
)

type WarehouseUseCase struct {
	repoPostgre WarehousePostgreRepo
}

func NewWarehouseUseCase(repoPostgre WarehousePostgreRepo) *WarehouseUseCase {
	return &WarehouseUseCase{
		repoPostgre,
	}
}

func (u *WarehouseUseCase) CreateWarehouse(ctx context.Context, warehouse *entity.Warehouse) error {
	err := warehouse.GenerateWarehouseID()
	if err != nil {
		return fmt.Errorf("failed to generate warehouse id: %w", err)
	}

	// save to postgres
	if err := u.repoPostgre.Save(ctx, warehouse); err != nil {
		return fmt.Errorf("failed to save warehouse: %w", err)
	}
	return nil
}

func (u *WarehouseUseCase) UpdateWarehouse(ctx context.Context, warehouse *entity.Warehouse) error {
	return u.repoPostgre.Update(ctx, warehouse)
}

func (u *WarehouseUseCase) GetWarehouseByID(ctx context.Context, id uuid.UUID) (*entity.Warehouse, error) {
	return u.repoPostgre.GetByID(ctx, id)
}

func (u *WarehouseUseCase) GetAllWarehouses(ctx context.Context) ([]*entity.Warehouse, error) {
	return u.repoPostgre.GetAll(ctx)
}

func (u *WarehouseUseCase) GetMainIDWarehouse(ctx context.Context) (uuid.UUID, error) {
	return u.repoPostgre.GetMainID(ctx)
}

// return nearest warehouse of zipcodes
func (u *WarehouseUseCase) GetNearestWarehouse(ctx context.Context, zipCodes []string) (map[string]string, error) {
	idAndZipCodes, err := u.repoPostgre.GetAllIDAndZipCode(ctx)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, zipCode := range zipCodes {
		nearest, err := utils.FindNearestWarehouseByZipCode(zipCode, idAndZipCodes)
		if err != nil {
			return nil, err
		}

		result[zipCode] = nearest
	}

	return result, nil
}
