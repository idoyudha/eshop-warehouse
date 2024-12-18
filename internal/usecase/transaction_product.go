package usecase

type TransactionProductUseCase struct {
	repoRedis          WarehouseRankRedisRepo
	repoMovePostgre    StockMovementPostgreRepo
	repoProductPostgre WarehouseProductPostgreRepo
}

func NewTransactionProductUseCase(repoRedis WarehouseRankRedisRepo, repoMovePostgre StockMovementPostgreRepo, repoProductPostgre WarehouseProductPostgreRepo) *TransactionProductUseCase {
	return &TransactionProductUseCase{
		repoRedis,
		repoMovePostgre,
		repoProductPostgre,
	}
}
