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
	gomock "go.uber.org/mock/gomock"
)

var (
	mockWarehouseProducts = []*entity.WarehouseProduct{
		{
			ID:              uuid.New(),
			ProductID:       uuid.New(),
			WarehouseID:     uuid.New(),
			ProductQuantity: 10,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		{
			ID:              uuid.New(),
			ProductID:       uuid.New(),
			WarehouseID:     uuid.New(),
			ProductQuantity: 20,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	mockWarehouseAddresses = []*entity.WarehouseAddressAndProductQty{
		{
			WarehouseID:     uuid.New(),
			ZipCode:         "12345",
			ProductName:     "Product A",
			ProductQuantity: 10,
		},
	}
)

type TestWarehouseProduct struct {
	name string
	mock func()
	res  any
	err  error
}

func warehouseProduct(t *testing.T) (*usecase.WarehouseProductUseCase, *MockWarehouseProductPostgreRepo) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockWarehouseProductPostgreRepo(mockCtl)
	warehouseProduct := usecase.NewWarehouseProductUseCase(repo)

	return warehouseProduct, repo
}

func TestCreateWarehouseProduct(t *testing.T) {
	// t.Parallell()
	warehouseProduct, repo := warehouseProduct(t)

	input := &entity.WarehouseProduct{
		ProductID:       uuid.New(),
		WarehouseID:     uuid.New(),
		ProductQuantity: 10,
	}

	tests := []TestWarehouseProduct{
		{
			name: "success",
			mock: func() {
				repo.EXPECT().
					Save(context.Background(), gomock.Any()).
					Return(nil)
			},
			res: nil,
			err: nil,
		},
		{
			name: "error",
			mock: func() {
				repo.EXPECT().
					Save(context.Background(), gomock.Any()).
					Return(errInternalServerError)
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

			err := warehouseProduct.CreateWarehouseProduct(context.Background(), input)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestGetAllWarehouseProducts(t *testing.T) {
	// t.Parallell()
	warehouseProduct, repo := warehouseProduct(t)

	tests := []TestWarehouseProduct{
		{
			name: "success",
			mock: func() {
				repo.EXPECT().
					GetAll(context.Background()).
					Return(mockWarehouseProducts, nil)
			},
			res: mockWarehouseProducts,
			err: nil,
		},
		{
			name: "error",
			mock: func() {
				repo.EXPECT().
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

			res, err := warehouseProduct.GetAllWarehouseProducts(context.Background())

			assert.Equal(t, tc.err, err)
			if err == nil {
				assert.NotNil(t, res)
				assert.Equal(t, tc.res, res)
				assert.Equal(t, len(tc.res.([]*entity.WarehouseProduct)), len(res))
				for i := 0; i < len(res); i++ {
					assert.Equal(t, tc.res.([]*entity.WarehouseProduct)[i].ID, res[i].ID)
					assert.Equal(t, tc.res.([]*entity.WarehouseProduct)[i].ProductID, res[i].ProductID)
					assert.Equal(t, tc.res.([]*entity.WarehouseProduct)[i].WarehouseID, res[i].WarehouseID)
					assert.Equal(t, tc.res.([]*entity.WarehouseProduct)[i].ProductQuantity, res[i].ProductQuantity)
					assert.Equal(t, tc.res.([]*entity.WarehouseProduct)[i].CreatedAt, res[i].CreatedAt)
					assert.Equal(t, tc.res.([]*entity.WarehouseProduct)[i].UpdatedAt, res[i].UpdatedAt)
				}
			} else {
				assert.Nil(t, res)
			}
		})
	}
}

func TestGetNearestWarehouseZipCodeByProductID(t *testing.T) {
	// t.Parallell()
	warehouseProduct, repo := warehouseProduct(t)

	inputZipCode := "12345"
	productID := uuid.New()

	mockAddresses := []*entity.WarehouseAddressAndProductQty{
		{
			WarehouseID:     uuid.New(),
			ZipCode:         "54321",
			ProductName:     "Product A",
			ProductQuantity: 10,
		},
	}

	tests := []TestWarehouseProduct{
		{
			name: "success",
			mock: func() {
				repo.EXPECT().
					GetWarehouseIDZipCodeAndQtyByProductID(context.Background(), productID).
					Return(mockAddresses, nil)
			},
			res: stringPtr("54321"),
			err: nil,
		},
		{
			name: "error get warehouse data",
			mock: func() {
				repo.EXPECT().
					GetWarehouseIDZipCodeAndQtyByProductID(context.Background(), productID).
					Return(nil, errInternalServerError)
			},
			res: nil,
			err: errors.New("failed to get warehouse and zipcode data: internal server error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallell()

			tc.mock()

			res, err := warehouseProduct.GetNearestWarehouseZipCodeByProductID(context.Background(), inputZipCode, productID)

			if tc.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, *tc.res.(*string), *res)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
