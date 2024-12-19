package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/entity"
	"github.com/idoyudha/eshop-warehouse/pkg/postgresql"
)

type TransactionProductPostgresRepo struct {
	*postgresql.Postgres
}

func NewTransactionProductPostgreRepo(client *postgresql.Postgres) *TransactionProductPostgresRepo {
	return &TransactionProductPostgresRepo{
		client,
	}
}

const (
	// locks the rows with FOR UPDATE
	queryLockSourceProduct = `
		SELECT id 
		FROM warehouse_products 
		WHERE product_id = $1 
		AND warehouse_id = $2 
		AND deleted_at IS NULL 
		FOR UPDATE`

	queryLockDestProduct = `
		SELECT id 
		FROM warehouse_products 
		WHERE product_id = $1 
		AND warehouse_id = $2 
		AND deleted_at IS NULL 
		FOR UPDATE`

	// updates
	queryUpdateSourceQuantity = `
		UPDATE warehouse_products 
		SET product_quantity = product_quantity - $1, 
		    updated_at = $2 
		WHERE product_id = $3 
		AND warehouse_id = $4 
		AND deleted_at IS NULL`

	queryUpdateDestQuantity = `
		UPDATE warehouse_products 
		SET product_quantity = COALESCE(product_quantity, 0) + $1, 
		    updated_at = $2
		WHERE product_id = $3 
		AND warehouse_id = $4
		AND deleted_at IS NULL`

	queryInsertDestProduct = `
		INSERT INTO warehouse_products (
			id, 
			warehouse_id, 
			product_id, 
			product_name, 
			product_quantity,
			created_at, 
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	queryInsertWarehouseMovement = `
		INSERT INTO stock_movements (
			id, 
			product_id, 
			product_name, 
			quantity, 
			from_warehouse_id, 
			to_warehouse_id, 
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`
)

// handling transfer from warehouse to warehouse
func (r *TransactionProductPostgresRepo) TransferIn(ctx context.Context, stockMovement *entity.StockMovement) error {
	// begin transaction
	tx, err := r.Conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. lock source product row if exists
	if err = tx.QueryRowContext(ctx, queryLockSourceProduct,
		stockMovement.ProductID, stockMovement.FromWarehouseID,
	).Scan(&struct{}{}); err != nil {
		return fmt.Errorf("failed to lock source product: %w", err)
	}

	// 2. lock destination product row if exists
	var destExist bool
	err = tx.QueryRowContext(ctx, queryLockDestProduct,
		stockMovement.ProductID, stockMovement.ToWarehouseID,
	).Scan(&struct{}{})
	destExist = err != sql.ErrNoRows
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check or lock destination product: %w", err)
	}

	// 3. update source quantity
	_, err = tx.ExecContext(ctx, queryUpdateSourceQuantity,
		stockMovement.Quantity, stockMovement.CreatedAt, stockMovement.ProductID, stockMovement.FromWarehouseID)
	if err != nil {
		return fmt.Errorf("failed to update source quantity: %w", err)
	}

	// 4. handle destination product
	if destExist {
		// update destination quantity
		_, err = tx.ExecContext(ctx, queryUpdateDestQuantity,
			stockMovement.Quantity, stockMovement.CreatedAt, stockMovement.ProductID, stockMovement.ToWarehouseID)
		if err != nil {
			return fmt.Errorf("failed to update destination quantity: %w", err)
		}
	} else {
		// create new product in destination warehouse
		newID, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("failed to generate uuid: %w", err)
		}
		_, err = tx.ExecContext(ctx, queryInsertDestProduct,
			newID,
			stockMovement.ToWarehouseID,
			stockMovement.ProductID,
			stockMovement.ProductName,
			stockMovement.Quantity,
			stockMovement.CreatedAt,
			stockMovement.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert destination product: %w", err)
		}
	}

	// 5. insert stock movement
	_, err = tx.ExecContext(ctx, queryInsertWarehouseMovement,
		stockMovement.ID,
		stockMovement.ProductID,
		stockMovement.ProductName,
		stockMovement.Quantity,
		stockMovement.FromWarehouseID,
		stockMovement.ToWarehouseID,
		stockMovement.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert stock movement: %w", err)
	}

	// commit transaction
	if errCommit := tx.Commit(); errCommit != nil {
		return fmt.Errorf("failed to commit transaction: %w", errCommit)
	}

	return nil
}

const queryInsertUserMovement = `
	INSERT INTO stock_movements (
		id, 
		product_id, 
		product_name, 
		quantity, 
		from_warehouse_id, 
		to_user_id,
		created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7)`

// handling transfer from warehouse to user
// if one warehouse is not enough products, then take it from another warehouse
// it will be multiple stock movement transactions
func (r *TransactionProductPostgresRepo) TransferOut(ctx context.Context, stockMovement []*entity.StockMovement) error {
	// begin transaction
	tx, err := r.Conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, movement := range stockMovement {
		// 1. lock source product row
		if err = tx.QueryRowContext(ctx, queryLockSourceProduct,
			movement.ProductID, movement.FromWarehouseID,
		).Scan(&struct{}{}); err != nil {
			return fmt.Errorf("failed to lock source product: %w", err)
		}

		// 2. update source quantity
		_, err = tx.ExecContext(ctx, queryUpdateSourceQuantity,
			movement.Quantity, movement.CreatedAt, movement.ProductID, movement.FromWarehouseID)
		if err != nil {
			return fmt.Errorf("failed to update source quantity: %w", err)
		}

		// 3. insert stock movement
		_, err = tx.ExecContext(ctx, queryInsertUserMovement,
			movement.ID,
			movement.ProductID,
			movement.ProductName,
			movement.Quantity,
			movement.FromWarehouseID,
			movement.ToUserID,
			movement.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert stock movement: %w", err)
		}
	}

	// commit transaction
	if errCommit := tx.Commit(); errCommit != nil {
		return fmt.Errorf("failed to commit transaction: %w", errCommit)
	}

	return nil
}
