package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock WarehouseProduct usecase
type mockWarehouseProductUsecase struct {
	mock.Mock
}

func (m *mockWarehouseProductUsecase) CreateWarehouseProduct(ctx context.Context, warehouseProduct *entity.WarehouseProduct) error {
	args := m.Called(ctx, warehouseProduct)
	return args.Error(0)
}

func (m *mockWarehouseProductUsecase) UpdateWarehouseProduct(ctx context.Context, warehouseProduct *entity.WarehouseProduct) error {
	args := m.Called(ctx, warehouseProduct)
	return args.Error(0)
}

func (m *mockWarehouseProductUsecase) UpdateWarehouseProductQuantity(ctx context.Context, warehouseProduct *entity.WarehouseProduct) error {
	args := m.Called(ctx, warehouseProduct)
	return args.Error(0)
}

func (m *mockWarehouseProductUsecase) GetAllWarehouseProducts(ctx context.Context) ([]*entity.WarehouseProduct, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.WarehouseProduct), args.Error(1)
}

func (m *mockWarehouseProductUsecase) GetWarehouseProductByProductID(ctx context.Context, productID uuid.UUID) ([]*entity.WarehouseProduct, error) {
	args := m.Called(ctx, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.WarehouseProduct), args.Error(1)
}

func (m *mockWarehouseProductUsecase) GetWarehouseProductByWarehouseID(ctx context.Context, warehouseID uuid.UUID) ([]*entity.WarehouseProduct, error) {
	args := m.Called(ctx, warehouseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.WarehouseProduct), args.Error(1)
}

func (m *mockWarehouseProductUsecase) GetWarehouseProductByProductIDAndWarehouseID(ctx context.Context, productID uuid.UUID, warehouseID uuid.UUID) (*entity.WarehouseProduct, error) {
	args := m.Called(ctx, productID, warehouseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.WarehouseProduct), args.Error(1)
}

func (m *mockWarehouseProductUsecase) GetNearestWarehouseZipCodeByProductID(ctx context.Context, zipCode string, productID uuid.UUID) (*string, error) {
	args := m.Called(ctx, zipCode, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func TestGetWarehouseProductByProductIDAndWarehouseID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	productID := uuid.New()
	warehouseID := uuid.New()
	mockProduct := &entity.WarehouseProduct{
		ID:              uuid.New(),
		ProductID:       productID,
		WarehouseID:     warehouseID,
		ProductQuantity: 4,
	}

	tests := []struct {
		name         string
		productID    string
		warehouseID  string
		expectedCode int
		setupMock    func(*mockWarehouseProductUsecase, *MockLogger)
	}{
		{
			name:         "success",
			productID:    productID.String(),
			warehouseID:  warehouseID.String(),
			expectedCode: http.StatusOK,
			setupMock: func(m *mockWarehouseProductUsecase, l *MockLogger) {
				m.On("GetWarehouseProductByProductIDAndWarehouseID",
					mock.Anything,
					productID,
					warehouseID,
				).Return(mockProduct, nil)
			},
		},
		{
			name:         "invalid product id",
			productID:    "invalid-uuid",
			warehouseID:  warehouseID.String(),
			expectedCode: http.StatusBadRequest,
			setupMock: func(m *mockWarehouseProductUsecase, l *MockLogger) {
				l.On("Error",
					mock.Anything,
					mock.Anything,
				).Return()
			},
		},
		{
			name:         "invalid warehouse id",
			productID:    productID.String(),
			warehouseID:  "invalid-uuid",
			expectedCode: http.StatusBadRequest,
			setupMock: func(m *mockWarehouseProductUsecase, l *MockLogger) {
				l.On("Error",
					mock.Anything,
					mock.Anything,
				).Return()
			},
		},
		{
			name:         "product not found",
			productID:    productID.String(),
			warehouseID:  warehouseID.String(),
			expectedCode: http.StatusInternalServerError,
			setupMock: func(m *mockWarehouseProductUsecase, l *MockLogger) {
				m.On("GetWarehouseProductByProductIDAndWarehouseID",
					mock.Anything,
					productID,
					warehouseID,
				).Return(nil, fmt.Errorf("product not found"))

				l.On("Error",
					mock.Anything,
					mock.Anything,
				).Return()
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockUC := new(mockWarehouseProductUsecase)
			mockLogger := NewMockLogger(t)

			tt.setupMock(mockUC, mockLogger)

			router := gin.New()
			handler := router.Group("/api/v1")
			newWarehouseProductRoutes(
				handler,
				mockUC,
				mockLogger,
				func(c *gin.Context) { c.Next() },
			)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/api/v1/warehouse-products/product/%s/warehouse/%s",
					tt.productID,
					tt.warehouseID,
				),
				nil,
			)

			router.ServeHTTP(w, req)

			if t.Failed() {
				t.Logf("Response Status: %d", w.Code)
				t.Logf("Response Body: %s", w.Body.String())
			}

			assert.Equal(t, tt.expectedCode, w.Code)
			mockUC.AssertExpectations(t)
			mockLogger.AssertExpectations(t)

			if tt.expectedCode == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "data")
			}
		})
	}
}
