package repo

import (
	"context"

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

const queryInsertWarehouse = `INSERT INTO warehouse (id, name, street, city, state, zip_code, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`

func (r *WarehousePostgreRepo) Save(ctx context.Context, warehouse *entity.Warehouse) error {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryInsertWarehouse)
	if errStmt != nil {
		return errStmt
	}
	defer stmt.Close()

	_, saveErr := stmt.ExecContext(ctx, warehouse.ID, warehouse.Name, warehouse.Street, warehouse.City, warehouse.State, warehouse.ZipCode, warehouse.CreatedAt, warehouse.UpdatedAt)
	if saveErr != nil {
		return saveErr
	}

	return nil
}

const queryUpdateWarehouse = `UPDATE warehouse SET name = $1, street = $2, city = $3, state = $4, zip_code = $5, updated_at = $6 WHERE id = $7;`

func (r *WarehousePostgreRepo) Update(ctx context.Context, warehouse *entity.Warehouse) error {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryUpdateWarehouse)
	if errStmt != nil {
		return errStmt
	}
	defer stmt.Close()

	_, updateErr := stmt.ExecContext(ctx, warehouse.Name, warehouse.Street, warehouse.City, warehouse.State, warehouse.ZipCode, warehouse.UpdatedAt, warehouse.ID)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

const queryGetByID = `SELECT id, name, street, city, state, zip_code, created_at, updated_at FROM warehouse WHERE id = $1 AND deleted_at IS NULL;`

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
		&warehouse.CreatedAt,
		&warehouse.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &warehouse, nil
}

const queryGetAll = `SELECT id, name, street, city, state, zip_code, created_at, updated_at FROM warehouse WHERE deleted_at IS NULL;`

func (r *WarehousePostgreRepo) GetAll(ctx context.Context) ([]*entity.Warehouse, error) {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryGetAll)
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
