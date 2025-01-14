package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/entity"
	"github.com/idoyudha/eshop-warehouse/internal/usecase"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

var (
	mockWarehouses = []*entity.Warehouse{
		{
			ID:        uuid.New(),
			Name:      "Warehouse A",
			Street:    "Street A",
			City:      "City A",
			State:     "State A",
			ZipCode:   "12345",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Name:      "Warehouse B",
			Street:    "Street B",
			City:      "City B",
			State:     "State B",
			ZipCode:   "123457",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
)

type TestWarehouse struct {
	name string
	mock func()
	res  any
	err  error
}

func warehouse(t *testing.T) (*usecase.WarehouseUseCase, *MockWarehouseRankRedisRepo, *MockWarehousePostgreRepo) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repoRedis := NewMockWarehouseRankRedisRepo(mockCtl)
	repoPostgre := NewMockWarehousePostgreRepo(mockCtl)
	warehouse := usecase.NewWarehouseUseCase(repoRedis, repoPostgre)

	return warehouse, repoRedis, repoPostgre
}

func TestGetAllWarehouses(t *testing.T) {
	// t.Parallell()
	warehouse, _, repoPostgre := warehouse(t)

	tests := []TestWarehouse{
		{
			name: "success",
			mock: func() {
				repoPostgre.EXPECT().
					GetAll(context.Background()).
					Return(mockWarehouses, nil)
			},
			res: mockWarehouses,
			err: nil,
		},
		{
			name: "error",
			mock: func() {
				repoPostgre.EXPECT().
					GetAll(context.Background()).
					Return(nil, errInternalServerError)
			},
			res: nil,
			err: errInternalServerError,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallell()

			tc.mock()

			res, err := warehouse.GetAllWarehouses(context.Background())

			assert.Equal(t, tc.err, err)
			if err == nil {
				assert.NotNil(t, res)
				assert.Equal(t, tc.res, res)
			} else {
				assert.Nil(t, res)
			}
		})
	}
}

func TestGetWarehouseByID(t *testing.T) {
	// t.Parallell()
	warehouse, _, repoPostgre := warehouse(t)

	tests := []TestWarehouse{
		{
			name: "success",
			mock: func() {
				repoPostgre.EXPECT().
					GetByID(context.Background(), mockWarehouses[0].ID).
					Return(mockWarehouses[0], nil)
			},
			res: mockWarehouses[0],
			err: nil,
		},
		{
			name: "error",
			mock: func() {
				repoPostgre.EXPECT().
					GetByID(context.Background(), mockWarehouses[0].ID).
					Return(nil, errInternalServerError)
			},
			res: nil,
			err: errInternalServerError,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallell()

			tc.mock()

			res, err := warehouse.GetWarehouseByID(context.Background(), mockWarehouses[0].ID)

			assert.Equal(t, tc.err, err)
			if err == nil {
				assert.NotNil(t, res)
				assert.Equal(t, tc.res, res)
			} else {
				assert.Nil(t, res)
			}
		})
	}
}
