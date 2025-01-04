package repo

import (
	"context"
	"log"

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

const queryInsertWarehouseProduct = `
	INSERT INTO warehouse_products (id, warehouse_id, product_id, product_sku, product_name, product_image_url, product_description, product_price, product_quantity, product_category_id, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);
`

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
		warehouseProduct.ProductSKU,
		warehouseProduct.ProductName,
		warehouseProduct.ProductImageURL,
		warehouseProduct.ProductDescription,
		warehouseProduct.ProductPrice,
		warehouseProduct.ProductQuantity,
		warehouseProduct.ProductCategoryID,
		warehouseProduct.CreatedAt,
		warehouseProduct.UpdatedAt,
	)
	if saveErr != nil {
		return saveErr
	}

	return nil
}

const queryUpdateNameAndPrice = `
	UPDATE warehouse_products 
	SET product_name = $1, product_image_url = $2, product_description = $3, product_price = $4, product_category_id = $5, updated_at = $6
	WHERE product_id = $7;
`

func (r *WarehouseProductPostgreRepo) Update(ctx context.Context, warehouseProduct *entity.WarehouseProduct) error {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryUpdateNameAndPrice)
	if errStmt != nil {
		return errStmt
	}
	defer stmt.Close()

	_, updateErr := stmt.ExecContext(ctx,
		warehouseProduct.ProductName,
		warehouseProduct.ProductImageURL,
		warehouseProduct.ProductDescription,
		warehouseProduct.ProductPrice,
		warehouseProduct.ProductCategoryID,
		warehouseProduct.UpdatedAt,
		warehouseProduct.ProductID,
	)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

const queryUpdateProductQuantity = `UPDATE warehouse_products SET product_quantity = $1, updated_at = $2 WHERE product_id = $3;`

func (r *WarehouseProductPostgreRepo) UpdateProductQuantity(ctx context.Context, warehouseProduct *entity.WarehouseProduct) error {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryUpdateProductQuantity)
	if errStmt != nil {
		return errStmt
	}
	defer stmt.Close()

	_, updateErr := stmt.ExecContext(ctx, warehouseProduct.ProductQuantity, warehouseProduct.UpdatedAt, warehouseProduct.ProductID)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

const queryGetAllWarehouseProducts = `
	SELECT id, warehouse_id, product_id, product_sku, product_name, product_image_url, product_description, product_price, product_quantity, product_category_id, created_at, updated_at
	FROM warehouse_products 
	WHERE deleted_at IS NULL;
`

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
			&warehouseProduct.ProductSKU,
			&warehouseProduct.ProductName,
			&warehouseProduct.ProductImageURL,
			&warehouseProduct.ProductDescription,
			&warehouseProduct.ProductPrice,
			&warehouseProduct.ProductQuantity,
			&warehouseProduct.ProductCategoryID,
			&warehouseProduct.CreatedAt,
			&warehouseProduct.UpdatedAt,
		)
		if err != nil {
			// TODO: handle error sql no rows
			return nil, err
		}
		warehouseProducts = append(warehouseProducts, &warehouseProduct)
	}

	return warehouseProducts, nil
}

const queryGetWarehouseProductByProductID = `
	SELECT id, warehouse_id, product_id, product_sku, product_name, product_image_url, product_description, product_price, product_quantity, product_category_id, created_at, updated_at
	FROM warehouse_products 
	WHERE product_id = $1 AND deleted_at IS NULL;
`

func (r *WarehouseProductPostgreRepo) GetByProductID(ctx context.Context, id uuid.UUID) ([]*entity.WarehouseProduct, error) {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryGetWarehouseProductByProductID)
	if errStmt != nil {
		return nil, errStmt
	}
	defer stmt.Close()

	var warehouseProducts []*entity.WarehouseProduct
	rows, err := stmt.QueryContext(ctx, id)
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
			&warehouseProduct.ProductSKU,
			&warehouseProduct.ProductName,
			&warehouseProduct.ProductImageURL,
			&warehouseProduct.ProductDescription,
			&warehouseProduct.ProductPrice,
			&warehouseProduct.ProductQuantity,
			&warehouseProduct.ProductCategoryID,
			&warehouseProduct.CreatedAt,
			&warehouseProduct.UpdatedAt,
		)
		if err != nil {
			// TODO: handle error sql no rows
			return nil, err
		}
		warehouseProducts = append(warehouseProducts, &warehouseProduct)
	}

	return warehouseProducts, nil
}

const queryGetWarehouseProductByWarehouseID = `
	SELECT id, warehouse_id, product_id, product_sku, product_name, product_image_url, product_description, product_price, product_quantity, product_category_id, created_at, updated_at
	FROM warehouse_products 
	WHERE warehouse_id = $1 AND deleted_at IS NULL;
`

func (r *WarehouseProductPostgreRepo) GetByWarehouseID(ctx context.Context, id uuid.UUID) ([]*entity.WarehouseProduct, error) {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryGetWarehouseProductByWarehouseID)
	if errStmt != nil {
		return nil, errStmt
	}
	defer stmt.Close()

	var warehouseProducts []*entity.WarehouseProduct
	rows, err := stmt.QueryContext(ctx, id)
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
			&warehouseProduct.ProductSKU,
			&warehouseProduct.ProductName,
			&warehouseProduct.ProductImageURL,
			&warehouseProduct.ProductDescription,
			&warehouseProduct.ProductPrice,
			&warehouseProduct.ProductQuantity,
			&warehouseProduct.ProductCategoryID,
			&warehouseProduct.CreatedAt,
			&warehouseProduct.UpdatedAt,
		)
		if err != nil {
			// TODO: handle error sql no rows
			return nil, err
		}
		warehouseProducts = append(warehouseProducts, &warehouseProduct)
	}

	return warehouseProducts, nil
}

const queryGetWarehouseProductByProductIDAndWarehouseID = `
	SELECT id, warehouse_id, product_id, product_sku, product_name, product_image_url, product_description, product_price, product_quantity, product_category_id, created_at, updated_at
	FROM warehouse_products 
	WHERE product_id = $1 AND warehouse_id = $2 AND deleted_at IS NULL;
`

func (r *WarehouseProductPostgreRepo) GetByProductIDAndWarehouseID(ctx context.Context, productID uuid.UUID, warehouseID uuid.UUID) (*entity.WarehouseProduct, error) {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryGetWarehouseProductByProductIDAndWarehouseID)
	if errStmt != nil {
		return nil, errStmt
	}
	defer stmt.Close()
	log.Printf("GetByProductIDAndWarehouseID productid => %s, warehouseid => %s", productID.String(), warehouseID.String())

	var warehouseProduct entity.WarehouseProduct
	err := stmt.QueryRowContext(ctx, productID, warehouseID).Scan(
		&warehouseProduct.ID,
		&warehouseProduct.WarehouseID,
		&warehouseProduct.ProductID,
		&warehouseProduct.ProductSKU,
		&warehouseProduct.ProductName,
		&warehouseProduct.ProductImageURL,
		&warehouseProduct.ProductDescription,
		&warehouseProduct.ProductPrice,
		&warehouseProduct.ProductQuantity,
		&warehouseProduct.ProductCategoryID,
		&warehouseProduct.CreatedAt,
		&warehouseProduct.UpdatedAt,
	)
	if err != nil {
		// TODO: handle error sql no rows
		return nil, err
	}

	return &warehouseProduct, nil
}

const queryGetWarehouseIDAndZipCodeByProductID = `
	SELECT warehouse_id, zip_code, product_name,product_quantity
	FROM warehouse_products
	JOIN warehouses
	ON warehouse_products.warehouse_id = warehouses.id
	WHERE warehouse_products.product_id = $1 AND warehouse_products.deleted_at IS NULL and warehouses.deleted_at IS NULL;
`

func (r *WarehouseProductPostgreRepo) GetWarehouseIDZipCodeAndQtyByProductID(ctx context.Context, productID uuid.UUID) ([]*entity.WarehouseAddressAndProductQty, error) {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryGetWarehouseIDAndZipCodeByProductID)
	if errStmt != nil {
		return nil, errStmt
	}
	defer stmt.Close()

	var warehouseAndProducts []*entity.WarehouseAddressAndProductQty
	rows, err := stmt.QueryContext(ctx, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var warehouseAndProduct entity.WarehouseAddressAndProductQty
		err := rows.Scan(
			&warehouseAndProduct.WarehouseID,
			&warehouseAndProduct.ZipCode,
			&warehouseAndProduct.ProductName,
			&warehouseAndProduct.ProductQuantity,
		)
		if err != nil {
			return nil, err
		}
		warehouseAndProducts = append(warehouseAndProducts, &warehouseAndProduct)
	}

	return warehouseAndProducts, nil
}

const queryGetTotalQuantityOfProductInAllWarehouse = `
	SELECT SUM(product_quantity) FROM warehouse_products WHERE product_id = $1 AND deleted_at IS NULL;
`

func (r *WarehouseProductPostgreRepo) GetTotalQuantityOfProductInAllWarehouse(ctx context.Context, productID uuid.UUID) (int, error) {
	stmt, errStmt := r.Conn.PrepareContext(ctx, queryGetTotalQuantityOfProductInAllWarehouse)
	if errStmt != nil {
		return 0, errStmt
	}
	defer stmt.Close()

	var totalQuantity int
	err := stmt.QueryRowContext(ctx, productID).Scan(&totalQuantity)
	if err != nil {
		return 0, err
	}

	return totalQuantity, nil
}
