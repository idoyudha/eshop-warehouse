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
	errMock = errors.New("mock error")

	mockID        = uuid.New()
	mockProductID = uuid.New()
	mockFromWHID  = uuid.New()
	mockToWHID    = uuid.New()
	mockToUserID  = uuid.New()
	mockTime      = time.Now()

	mockStockMovements = []*entity.StockMovement{
		{
			ID:              mockID,
			ProductID:       mockProductID,
			ProductName:     "Test Product",
			Quantity:        10,
			FromWarehouseID: mockFromWHID,
			ToWarehouseID:   mockToWHID,
			ToUserID:        mockToUserID,
			CreatedAt:       mockTime,
		},
	}
)

type TestStockMovement struct {
	name string
	mock func()
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
	stockMovement, repo := stockMovement(t)

	tests := []TestStockMovement{
		{
			name: "success",
			mock: func() {
				repo.EXPECT().
					GetAll(gomock.Any()).
					Return(mockStockMovements, nil)
			},
			res: mockStockMovements,
			err: nil,
		},
		{
			name: "error",
			mock: func() {
				repo.EXPECT().
					GetAll(gomock.Any()).
					Return(nil, errMock)
			},
			res: nil,
			err: errMock,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			tc.mock()
			res, err := stockMovement.GetAllStockMovements(context.Background())

			assert.Equal(t, tc.err, err)
			if err == nil {
				assert.NotNil(t, res)
				assert.Equal(t, tc.res, res)
				assert.Equal(t, len(tc.res), len(res))
				if len(res) > 0 {
					assert.Equal(t, tc.res[0].ID, res[0].ID)
					assert.Equal(t, tc.res[0].ProductID, res[0].ProductID)
					assert.Equal(t, tc.res[0].ProductName, res[0].ProductName)
					assert.Equal(t, tc.res[0].Quantity, res[0].Quantity)
					assert.Equal(t, tc.res[0].FromWarehouseID, res[0].FromWarehouseID)
					assert.Equal(t, tc.res[0].ToWarehouseID, res[0].ToWarehouseID)
					assert.Equal(t, tc.res[0].ToUserID, res[0].ToUserID)
					assert.Equal(t, tc.res[0].CreatedAt, res[0].CreatedAt)
				}
			} else {
				assert.Nil(t, res)
			}
		})
	}
}
