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

func NewWarehousePostgreRepo(client *postgresql.Postgres) *WarehousePostgreRepo {
	return &WarehousePostgreRepo{
		client,
	}
}

const queryInsertWarehouse = `
	INSERT INTO warehouses (id, name, street, city, state, zip_code, is_main_warehouse, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
`

func (r *WarehousePostgreRepo) Save(ctx context.Context, warehouse *entity.Warehouse) error {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryInsertWarehouse)
	if errStmt != nil {
		return fmt.Errorf("failed to prepare statement: %w", errStmt)
	}
	defer stmt.Close()

	_, saveErr := stmt.ExecContext(ctx,
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
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryUpdateWarehouse)
	if errStmt != nil {
		return errStmt
	}
	defer stmt.Close()

	_, updateErr := stmt.ExecContext(ctx, warehouse.Name, warehouse.Street, warehouse.UpdatedAt, warehouse.ID)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

const queryGetByID = `SELECT id, name, street, city, state, zip_code, is_main_warehouse, created_at, updated_at FROM warehouses WHERE id = $1 AND deleted_at IS NULL;`

func (r *WarehousePostgreRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Warehouse, error) {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryGetByID)
	if errStmt != nil {
		return nil, errStmt
	}
	defer stmt.Close()

	var warehouse entity.Warehouse
	err := stmt.QueryRowContext(ctx, id).Scan(
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
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryGetAllWarehouse)
	if errStmt != nil {
		return nil, errStmt
	}
	defer stmt.Close()

	var warehouses []*entity.Warehouse
	rows, err := stmt.QueryContext(ctx)
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
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryGetAllExceptMainWarehouse)
	if errStmt != nil {
		return nil, errStmt
	}
	defer stmt.Close()

	var warehouses []*entity.Warehouse
	rows, err := stmt.QueryContext(ctx)
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
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryGetMainIDWarehouse)
	if errStmt != nil {
		return uuid.Nil, errStmt
	}
	defer stmt.Close()

	var id uuid.UUID
	err := stmt.QueryRowContext(ctx).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

const queryGetAllWarehouseIDAndZipCode = `SELECT id, zip_code FROM warehouses WHERE deleted_at IS NULL ORDER BY zip_code ASC;`

func (r *WarehousePostgreRepo) GetAllIDAndZipCode(ctx context.Context) ([]*entity.Warehouse, error) {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryGetAllWarehouseIDAndZipCode)
	if errStmt != nil {
		return nil, errStmt
	}
	defer stmt.Close()

	var warehouses []*entity.Warehouse
	rows, err := stmt.QueryContext(ctx)
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
