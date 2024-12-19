package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/entity"
	"github.com/idoyudha/eshop-warehouse/pkg/postgresql"
)

type WarehouseProductPostgreRepo struct {
	*postgresql.Postgres
}

func NewWarehouseProductPostgreRepo(client *postgresql.Postgres) *WarehouseProductPostgreRepo {
	return &WarehouseProductPostgreRepo{
		client,
	}
}

const queryInsertWarehouseProduct = `INSERT INTO warehouse_products (id, warehouse_id, product_id, product_name, product_price, product_quantity, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`

func (r *WarehouseProductPostgreRepo) Save(ctx context.Context, warehouseProduct *entity.WarehouseProduct) error {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryInsertWarehouseProduct)
	if errStmt != nil {
		return errStmt
	}
	defer stmt.Close()

	_, saveErr := stmt.ExecContext(ctx,
		warehouseProduct.ID,
		warehouseProduct.WarehouseID,
		warehouseProduct.ProductID,
		warehouseProduct.ProductName,
		warehouseProduct.ProductPrice,
		warehouseProduct.ProductQuantity,
		warehouseProduct.CreatedAt,
		warehouseProduct.UpdatedAt,
	)
	if saveErr != nil {
		return saveErr
	}

	return nil
}

const queryUpdateProductQuantity = `UPDATE warehouse_products SET product_quantity = $1, updated_at = $2 WHERE warehouse_id = $3;`

func (r *WarehouseProductPostgreRepo) UpdateProductQuantity(ctx context.Context, warehouseProduct *entity.WarehouseProduct) error {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryUpdateProductQuantity)
	if errStmt != nil {
		return errStmt
	}
	defer stmt.Close()

	_, updateErr := stmt.ExecContext(ctx, warehouseProduct.ProductQuantity, warehouseProduct.UpdatedAt, warehouseProduct.WarehouseID)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

const queryGetAllWarehouseProducts = `SELECT id, warehouse_id, product_id, product_name, product_price, product_quantity, created_at, updated_at FROM warehouse_products WHERE deleted_at IS NULL;`

func (r *WarehouseProductPostgreRepo) GetAll(ctx context.Context) ([]*entity.WarehouseProduct, error) {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryGetAllWarehouseProducts)
	if errStmt != nil {
		return nil, errStmt
	}
	defer stmt.Close()

	var warehouseProducts []*entity.WarehouseProduct
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var warehouseProduct entity.WarehouseProduct
		err := rows.Scan(
			&warehouseProduct.ID,
			&warehouseProduct.WarehouseID,
			&warehouseProduct.ProductID,
			&warehouseProduct.ProductName,
			&warehouseProduct.ProductPrice,
			&warehouseProduct.ProductQuantity,
			&warehouseProduct.CreatedAt,
			&warehouseProduct.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		warehouseProducts = append(warehouseProducts, &warehouseProduct)
	}

	return warehouseProducts, nil
}

const queryGetWarehouseProductByProductID = `SELECT id, warehouse_id, product_id, product_name, product_price, product_quantity, created_at, updated_at FROM warehouse_products WHERE product_id = $1 AND deleted_at IS NULL;`

func (r *WarehouseProductPostgreRepo) GetByProductID(ctx context.Context, id uuid.UUID) ([]*entity.WarehouseProduct, error) {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryGetWarehouseProductByProductID)
	if errStmt != nil {
		return nil, errStmt
	}
	defer stmt.Close()

	var warehouseProducts []*entity.WarehouseProduct
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var warehouseProduct entity.WarehouseProduct
		err := rows.Scan(
			&warehouseProduct.ID,
			&warehouseProduct.WarehouseID,
			&warehouseProduct.ProductID,
			&warehouseProduct.ProductName,
			&warehouseProduct.ProductPrice,
			&warehouseProduct.ProductQuantity,
			&warehouseProduct.CreatedAt,
			&warehouseProduct.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		warehouseProducts = append(warehouseProducts, &warehouseProduct)
	}

	return warehouseProducts, nil
}

const queryGetWarehouseProductByWarehouseID = `SELECT id, warehouse_id, product_id, product_name, product_price, product_quantity, created_at, updated_at FROM warehouse_products WHERE warehouse_id = $1 AND deleted_at IS NULL;`

func (r *WarehouseProductPostgreRepo) GetByWarehouseID(ctx context.Context, id uuid.UUID) ([]*entity.WarehouseProduct, error) {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryGetWarehouseProductByWarehouseID)
	if errStmt != nil {
		return nil, errStmt
	}
	defer stmt.Close()

	var warehouseProducts []*entity.WarehouseProduct
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var warehouseProduct entity.WarehouseProduct
		err := rows.Scan(
			&warehouseProduct.ID,
			&warehouseProduct.WarehouseID,
			&warehouseProduct.ProductID,
			&warehouseProduct.ProductName,
			&warehouseProduct.ProductPrice,
			&warehouseProduct.ProductQuantity,
			&warehouseProduct.CreatedAt,
			&warehouseProduct.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		warehouseProducts = append(warehouseProducts, &warehouseProduct)
	}

	return warehouseProducts, nil
}

const queryGetWarehouseProductByProductIDAndWarehouseID = `SELECT id, warehouse_id, product_id, product_name, product_price, product_quantity, created_at, updated_at FROM warehouse_products WHERE product_id = $1 AND warehouse_id = $2 AND deleted_at IS NULL;`

func (r *WarehouseProductPostgreRepo) GetByProductIDAndWarehouseID(ctx context.Context, productID uuid.UUID, warehouseID uuid.UUID) (*entity.WarehouseProduct, error) {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryGetWarehouseProductByProductIDAndWarehouseID)
	if errStmt != nil {
		return nil, errStmt
	}
	defer stmt.Close()

	var warehouseProduct entity.WarehouseProduct
	err := stmt.QueryRowContext(ctx, productID, warehouseID).Scan(
		&warehouseProduct.ID,
		&warehouseProduct.WarehouseID,
		&warehouseProduct.ProductID,
		&warehouseProduct.ProductName,
		&warehouseProduct.ProductPrice,
		&warehouseProduct.ProductQuantity,
		&warehouseProduct.CreatedAt,
		&warehouseProduct.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &warehouseProduct, nil
}
