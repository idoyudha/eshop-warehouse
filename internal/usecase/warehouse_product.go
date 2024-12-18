package usecase

type WarehouseProductUseCase struct {
	repoPostgre WarehouseProductPostgreRepo
}

func NewWarehouseProductUseCase(repoPostgre WarehouseProductPostgreRepo) *WarehouseProductUseCase {
	return &WarehouseProductUseCase{
		repoPostgre,
	}
}
