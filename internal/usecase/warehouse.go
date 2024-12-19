package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/entity"
	"github.com/idoyudha/eshop-warehouse/internal/utils"
)

type WarehouseUseCase struct {
	repoRedis   WarehouseRankRedisRepo
	repoPostgre WarehousePostgreRepo
}

func NewWarehouseUseCase(repoRedis WarehouseRankRedisRepo, repoPostgre WarehousePostgreRepo) *WarehouseUseCase {
	return &WarehouseUseCase{
		repoRedis,
		repoPostgre,
	}
}

func (u *WarehouseUseCase) CreateWarehouse(ctx context.Context, warehouse *entity.Warehouse) error {
	// check if this is the first warehouse
	existingWarehouses, err := u.repoPostgre.GetAllExceptMain(ctx)
	if err != nil {
		return err
	}

	// save to postgres
	if err := u.repoPostgre.Save(ctx, warehouse); err != nil {
		return err
	}

	// if this is the first warehouse
	if len(existingWarehouses) == 0 {
		distances := map[string]float64{
			warehouse.ZipCode: 0,
		}

		err = u.repoRedis.UpdateWarehouseRanks(ctx, warehouse.ZipCode, warehouse.ID.String(), distances)
		if err != nil {
			// TODO: handled save record to the postgres
			// cancel the save or try use retry or send event to kafka
			return fmt.Errorf("failed to set initial warehouse rank: %w", err)
		}
		return nil
	}

	zipCodes := make(map[string]string)
	for _, w := range existingWarehouses {
		zipCodes[w.ZipCode] = w.ZipCode
	}

	distances, err := utils.CalculateZipCodeDistance(warehouse.ZipCode, zipCodes)
	if err != nil {
		return fmt.Errorf("failed to calculate distances: %w", err)
	}

	// update new warehouse rankings
	err = u.repoRedis.UpdateWarehouseRanks(ctx, warehouse.ZipCode, warehouse.ID.String(), distances)
	if err != nil {
		// TODO: handled save record to the postgres
		// cancel the save or try use retry or send event to kafka
		return fmt.Errorf("failed to set initial warehouse rank: %w", err)
	}

	// update rankings for existing warehouses
	for _, existingWarehouse := range existingWarehouses {
		// Calculate distance from existing warehouse to new warehouse
		distancesFromExisting, err := utils.CalculateZipCodeDistance(
			existingWarehouse.ZipCode,
			map[string]string{warehouse.ZipCode: warehouse.ZipCode},
		)
		if err != nil {
			return fmt.Errorf("failed to calculate reverse distances: %w", err)
		}

		// Update Redis rankings for existing warehouse
		err = u.repoRedis.UpdateWarehouseRanks(
			ctx,
			existingWarehouse.ZipCode,
			existingWarehouse.ID.String(),
			distancesFromExisting,
		)
		if err != nil {
			return fmt.Errorf("failed to update reverse warehouse ranks: %w", err)
		}
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
