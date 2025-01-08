package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/entity"
	"github.com/idoyudha/eshop-warehouse/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var (
	errInternalServerError = errors.New("internal server error")

	mockStockMovements = []*entity.StockMovement{
		{
			ID:              uuid.New(),
			ProductID:       uuid.New(),
			ProductName:     "Product A",
			Quantity:        10,
			FromWarehouseID: uuid.New(),
			ToWarehouseID:   uuid.New(),
			ToUserID:        uuid.New(),
			CreatedAt:       time.Now(),
		},
		{
			ID:              uuid.New(),
			ProductID:       uuid.New(),
			ProductName:     "Product B",
			Quantity:        20,
			FromWarehouseID: uuid.New(),
			ToWarehouseID:   uuid.New(),
			ToUserID:        uuid.New(),
			CreatedAt:       time.Now(),
		},
	}
)

type TestStockMovement struct {
	name string
	mock func(*MockStockMovementPostgreRepo)
	res  []*entity.StockMovement
	err  error
}

func stockMovement(t *testing.T) (*usecase.StockMovementUseCase, *MockStockMovementPostgreRepo) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockStockMovementPostgreRepo(mockCtl)
	stockMovement := usecase.NewStockMovementUseCase(repo)

	return stockMovement, repo
}

func TestGetAllStockMovements(t *testing.T) {
	// allow this function run in parallel with other test function
	t.Parallel()

	tests := []TestStockMovement{
		{
			name: "success",
			mock: func(repo *MockStockMovementPostgreRepo) {
				repo.EXPECT().
					GetAll(gomock.Any()).
					Return(mockStockMovements, nil)
			},
			res: mockStockMovements,
			err: nil,
		},
		{
			name: "error",
			mock: func(repo *MockStockMovementPostgreRepo) {
				repo.EXPECT().
					GetAll(gomock.Any()).
					Return(nil, errInternalServerError)
			},
			res: nil,
			err: errInternalServerError,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// test case will run in parallel
			t.Parallel()
			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()

			repo := NewMockStockMovementPostgreRepo(mockCtl)
			stockMovement := usecase.NewStockMovementUseCase(repo)

			tc.mock(repo)
			res, err := stockMovement.GetAllStockMovements(context.Background())

			assert.Equal(t, tc.err, err)
			if err == nil {
				assert.NotNil(t, res)
				assert.Equal(t, tc.res, res)
				assert.Equal(t, len(tc.res), len(res))
				for i := 0; i < len(tc.res); i++ {
					assert.Equal(t, tc.res[i].ID, res[i].ID)
					assert.Equal(t, tc.res[i].ProductID, res[i].ProductID)
					assert.Equal(t, tc.res[i].ProductName, res[i].ProductName)
					assert.Equal(t, tc.res[i].Quantity, res[i].Quantity)
					assert.Equal(t, tc.res[i].FromWarehouseID, res[i].FromWarehouseID)
					assert.Equal(t, tc.res[i].ToWarehouseID, res[i].ToWarehouseID)
					assert.Equal(t, tc.res[i].ToUserID, res[i].ToUserID)
					assert.Equal(t, tc.res[i].CreatedAt, res[i].CreatedAt)
				}
			} else {
				assert.Nil(t, res)
			}
		})
	}
}

// func TestGetByProductID(t *testing.T) {
// 	stockMovement, repo := stockMovement(t)
// }
