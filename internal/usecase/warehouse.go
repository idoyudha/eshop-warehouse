package usecase

type WarehouseUseCase struct {
	repoRedis   WarehouseRankRedisRepo
	repoPostgre WarehousePostgreRepo
}

func NewWarehouseUseCase(repoRedis WarehouseRankRedisRepo, repoPostgre WarehousePostgreRepo) *WarehouseUseCase {
	return &WarehouseUseCase{
		repoRedis,
		repoPostgre,
	}
}
