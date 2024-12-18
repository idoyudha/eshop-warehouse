package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/entity"
)

type StockMovementUseCase struct {
	repoMovePostgre StockMovementPostgreRepo
}

func NewStockMovementUseCase(repoMovePostgre StockMovementPostgreRepo) *StockMovementUseCase {
	return &StockMovementUseCase{
		repoMovePostgre,
	}
}

func (u *StockMovementUseCase) GetAllStockMovements(ctx context.Context) ([]*entity.StockMovement, error) {
	return u.repoMovePostgre.GetAll(ctx)
}

func (u *StockMovementUseCase) GetStockMovementsByProductID(ctx context.Context, productID uuid.UUID) ([]*entity.StockMovement, error) {
	return u.repoMovePostgre.GetByProductID(ctx, productID)
}

func (u *StockMovementUseCase) GetStockMovementsBySourceID(ctx context.Context, sourceID uuid.UUID) ([]*entity.StockMovement, error) {
	return u.repoMovePostgre.GetBySourceID(ctx, sourceID)
}

func (u *StockMovementUseCase) GetStockMovementsByDestinationID(ctx context.Context, destinationID uuid.UUID) ([]*entity.StockMovement, error) {
	return u.repoMovePostgre.GetByDestinationID(ctx, destinationID)
}
