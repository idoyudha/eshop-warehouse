package usecase

type StockMovementUseCase struct {
	repoRedis          WarehouseRankRedisRepo
	repoMovePostgre    StockMovementPostgreRepo
	repoProductPostgre WarehouseProductPostgreRepo
}

func NewStockMovementUseCase(repoRedis WarehouseRankRedisRepo, repoMovePostgre StockMovementPostgreRepo, repoProductPostgre WarehouseProductPostgreRepo) *StockMovementUseCase {
	return &StockMovementUseCase{
		repoRedis,
		repoMovePostgre,
		repoProductPostgre,
	}
}
