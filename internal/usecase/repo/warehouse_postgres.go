package repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/entity"
	"github.com/idoyudha/eshop-warehouse/pkg/postgresql"
)

type WarehousePostgreRepo struct {
	*postgresql.Postgres
}

func NewWarehousePostgreRepo(pg *postgresql.Postgres) *WarehousePostgreRepo {
	return &WarehousePostgreRepo{
		pg,
	}
}

const queryInsertWarehouse = `
	INSERT INTO warehouses (id, name, street, city, state, zip_code, is_main_warehouse, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
`

func (r *WarehousePostgreRepo) Save(ctx context.Context, warehouse *entity.Warehouse) error {
	_, saveErr := r.Pool.Exec(ctx,
		queryInsertWarehouse,
		warehouse.ID,
		warehouse.Name,
		warehouse.Street,
		warehouse.City,
		warehouse.State,
		warehouse.ZipCode,
		warehouse.IsMainWarehouse,
		warehouse.CreatedAt,
		warehouse.UpdatedAt,
	)
	if saveErr != nil {
		return fmt.Errorf("failed to save warehouse: %w", saveErr)
	}

	return nil
}

const queryUpdateWarehouse = `UPDATE warehouses SET name = $1, street = $2, updated_at = $3 WHERE id = $4;`

func (r *WarehousePostgreRepo) Update(ctx context.Context, warehouse *entity.Warehouse) error {
	_, updateErr := r.Pool.Exec(ctx, queryUpdateWarehouse,
		warehouse.Name, warehouse.Street, warehouse.UpdatedAt, warehouse.ID)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

const queryGetByID = `SELECT id, name, street, city, state, zip_code, is_main_warehouse, created_at, updated_at FROM warehouses WHERE id = $1 AND deleted_at IS NULL;`

func (r *WarehousePostgreRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Warehouse, error) {
	var warehouse entity.Warehouse
	err := r.Pool.QueryRow(ctx, queryGetByID, id).Scan(
		&warehouse.ID,
		&warehouse.Name,
		&warehouse.Street,
		&warehouse.City,
		&warehouse.State,
		&warehouse.ZipCode,
		&warehouse.IsMainWarehouse,
		&warehouse.CreatedAt,
		&warehouse.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &warehouse, nil
}

const queryGetAllWarehouse = `SELECT id, name, street, city, state, zip_code, is_main_warehouse, created_at, updated_at FROM warehouses WHERE deleted_at IS NULL;`

func (r *WarehousePostgreRepo) GetAll(ctx context.Context) ([]*entity.Warehouse, error) {
	var warehouses []*entity.Warehouse
	rows, err := r.Pool.Query(ctx, queryGetAllWarehouse)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var warehouse entity.Warehouse
		err := rows.Scan(
			&warehouse.ID,
			&warehouse.Name,
			&warehouse.Street,
			&warehouse.City,
			&warehouse.State,
			&warehouse.ZipCode,
			&warehouse.IsMainWarehouse,
			&warehouse.CreatedAt,
			&warehouse.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		warehouses = append(warehouses, &warehouse)
	}

	return warehouses, nil
}

const queryGetAllExceptMainWarehouse = `SELECT id, name, street, city, state, zip_code, is_main_warehouse, created_at, updated_at FROM warehouses WHERE is_main_warehouse = false AND deleted_at IS NULL;`

func (r *WarehousePostgreRepo) GetAllExceptMain(ctx context.Context) ([]*entity.Warehouse, error) {
	var warehouses []*entity.Warehouse
	rows, err := r.Pool.Query(ctx, queryGetAllExceptMainWarehouse)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var warehouse entity.Warehouse
		err := rows.Scan(
			&warehouse.ID,
			&warehouse.Name,
			&warehouse.Street,
			&warehouse.City,
			&warehouse.State,
			&warehouse.ZipCode,
			&warehouse.IsMainWarehouse,
			&warehouse.CreatedAt,
			&warehouse.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		warehouses = append(warehouses, &warehouse)
	}

	return warehouses, nil
}

const queryGetMainIDWarehouse = `SELECT id FROM warehouses WHERE is_main_warehouse = true AND deleted_at IS NULL;`

func (r *WarehousePostgreRepo) GetMainID(ctx context.Context) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.Pool.QueryRow(ctx, queryGetMainIDWarehouse).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

const queryGetAllWarehouseIDAndZipCode = `SELECT id, zip_code FROM warehouses WHERE deleted_at IS NULL ORDER BY zip_code ASC;`

func (r *WarehousePostgreRepo) GetAllIDAndZipCode(ctx context.Context) ([]*entity.Warehouse, error) {
	var warehouses []*entity.Warehouse
	rows, err := r.Pool.Query(ctx, queryGetAllWarehouseIDAndZipCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var warehouse entity.Warehouse
		err := rows.Scan(
			&warehouse.ID,
			&warehouse.ZipCode,
		)
		if err != nil {
			return nil, err
		}
		warehouses = append(warehouses, &warehouse)
	}

	return warehouses, nil
}
