package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/entity"
	"github.com/idoyudha/eshop-warehouse/pkg/postgresql"
)

type StockMovementPostgreRepo struct {
	*postgresql.Postgres
}

func NewStockMovementPostgreRepo(pg *postgresql.Postgres) *StockMovementPostgreRepo {
	return &StockMovementPostgreRepo{
		pg,
	}
}

const queryGetAllStockMovements = `SELECT * FROM stock_movements;`

func (r *StockMovementPostgreRepo) GetAll(ctx context.Context) ([]*entity.StockMovement, error) {
	rows, err := r.Pool.Query(ctx, queryGetAllStockMovements)
	if err != nil {
		return nil, err
	}
	rows.Close()

	var stockMovements []*entity.StockMovement
	for rows.Next() {
		var stockMovement entity.StockMovement
		if err := rows.Scan(
			&stockMovement.ID,
			&stockMovement.ProductID,
			&stockMovement.ProductName,
			&stockMovement.Quantity,
			&stockMovement.FromWarehouseID,
			&stockMovement.ToWarehouseID,
			&stockMovement.ToUserID,
			&stockMovement.CreatedAt,
		); err != nil {
			return nil, err
		}
		stockMovements = append(stockMovements, &stockMovement)
	}

	return stockMovements, nil
}

const queryGetByProductID = `SELECT * FROM stock_movements WHERE product_id = $1;`

func (r *StockMovementPostgreRepo) GetByProductID(ctx context.Context, productID uuid.UUID) ([]*entity.StockMovement, error) {
	rows, err := r.Pool.Query(ctx, queryGetByProductID, productID)
	if err != nil {
		return nil, err
	}
	rows.Close()

	var stockMovements []*entity.StockMovement
	for rows.Next() {
		var stockMovement entity.StockMovement
		if err := rows.Scan(
			&stockMovement.ID,
			&stockMovement.ProductID,
			&stockMovement.ProductName,
			&stockMovement.Quantity,
			&stockMovement.FromWarehouseID,
			&stockMovement.ToWarehouseID,
			&stockMovement.ToUserID,
			&stockMovement.CreatedAt,
		); err != nil {
			return nil, err
		}
		stockMovements = append(stockMovements, &stockMovement)
	}

	return stockMovements, nil
}

const queryGetBySourceID = `SELECT * FROM stock_movements WHERE from_warehouse_id = $1;`

func (r *StockMovementPostgreRepo) GetBySourceID(ctx context.Context, sourceID uuid.UUID) ([]*entity.StockMovement, error) {
	rows, err := r.Pool.Query(ctx, queryGetBySourceID, sourceID)
	if err != nil {
		return nil, err
	}
	rows.Close()

	var stockMovements []*entity.StockMovement
	for rows.Next() {
		var stockMovement entity.StockMovement
		if err := rows.Scan(
			&stockMovement.ID,
			&stockMovement.ProductID,
			&stockMovement.ProductName,
			&stockMovement.Quantity,
			&stockMovement.FromWarehouseID,
			&stockMovement.ToWarehouseID,
			&stockMovement.ToUserID,
			&stockMovement.CreatedAt,
		); err != nil {
			return nil, err
		}
		stockMovements = append(stockMovements, &stockMovement)
	}

	return stockMovements, nil
}

const queryGetByDestinationID = `SELECT * FROM stock_movements WHERE to_warehouse_id = $1;`

func (r *StockMovementPostgreRepo) GetByDestinationID(ctx context.Context, destinationID uuid.UUID) ([]*entity.StockMovement, error) {
	rows, err := r.Pool.Query(ctx, queryGetByDestinationID, destinationID)
	if err != nil {
		return nil, err
	}
	rows.Close()

	var stockMovements []*entity.StockMovement
	for rows.Next() {
		var stockMovement entity.StockMovement
		if err := rows.Scan(
			&stockMovement.ID,
			&stockMovement.ProductID,
			&stockMovement.ProductName,
			&stockMovement.Quantity,
			&stockMovement.FromWarehouseID,
			&stockMovement.ToWarehouseID,
			&stockMovement.ToUserID,
			&stockMovement.CreatedAt,
		); err != nil {
			return nil, err
		}
		stockMovements = append(stockMovements, &stockMovement)
	}

	return stockMovements, nil
}
