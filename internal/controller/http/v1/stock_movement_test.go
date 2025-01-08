package v1

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/entity"
	"github.com/idoyudha/eshop-warehouse/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockStockMovementUsecase struct {
	mock.Mock
}

func (m *mockStockMovementUsecase) GetAllStockMovements(ctx context.Context) ([]*entity.StockMovement, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*entity.StockMovement), args.Error(1)
}

func (m *mockStockMovementUsecase) GetStockMovementsByProductID(ctx context.Context, productID uuid.UUID) ([]*entity.StockMovement, error) {
	args := m.Called(ctx, productID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*entity.StockMovement), args.Error(1)
}

func (m *mockStockMovementUsecase) GetStockMovementsBySourceID(ctx context.Context, sourceID uuid.UUID) ([]*entity.StockMovement, error) {
	args := m.Called(ctx, sourceID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*entity.StockMovement), args.Error(1)
}

func (m *mockStockMovementUsecase) GetStockMovementsByDestinationID(ctx context.Context, destinationID uuid.UUID) ([]*entity.StockMovement, error) {
	args := m.Called(ctx, destinationID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*entity.StockMovement), args.Error(1)
}

// interface implementation
var _ usecase.StockMovement = (*mockStockMovementUsecase)(nil)

type mockTransactionProductUsecase struct {
	mock.Mock
}

func (m *mockTransactionProductUsecase) MoveIn(ctx context.Context, stockMovement *entity.StockMovement) error {
	args := m.Called(ctx, stockMovement)
	return args.Error(0)
}

func (m *mockTransactionProductUsecase) MoveOut(ctx context.Context, stockMovements []*entity.StockMovement, zipCode string) error {
	args := m.Called(ctx, stockMovements, zipCode)
	return args.Error(0)
}

type testStockMovement struct {
	name             string
	inputJSON        string
	expectedStatus   int
	expectedResponse interface{}
	mockBehavior     func(*mockTransactionProductUsecase)
}

var _ usecase.TransactionProduct = (*mockTransactionProductUsecase)(nil)

func TestCreateStockMovementIn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		inputJSON      string
		expectedStatus int
		mockBehavior   func(*mockTransactionProductUsecase, *MockLogger)
	}{
		{
			name: "Success",
			inputJSON: `{
                "product_id": "019444a2-e318-79b5-8fe4-b32716306083",
                "product_name": "Product A",
                "quantity": 10,
                "from_warehouse_id": "019444a3-6a3f-7249-b694-f6f071d8eb79",
                "to_warehouse_id": "019444a3-a5dc-7e93-bcc3-fec46dddd299"
            }`,
			expectedStatus: http.StatusCreated,
			mockBehavior: func(m *mockTransactionProductUsecase, l *MockLogger) {
				m.On("MoveIn",
					mock.Anything,
					mock.MatchedBy(func(sm *entity.StockMovement) bool {
						return sm.ProductName == "Product A" &&
							sm.ProductID.String() == "019444a2-e318-79b5-8fe4-b32716306083" &&
							sm.Quantity == 10 &&
							sm.FromWarehouseID.String() == "019444a3-6a3f-7249-b694-f6f071d8eb79" &&
							sm.ToWarehouseID.String() == "019444a3-a5dc-7e93-bcc3-fec46dddd299"
					}),
				).Return(nil)
			},
		},
		{
			name: "Invalid JSON",
			inputJSON: `{
                "product_id": "invalid-uuid",
                "quantity": "invalid"
            }`,
			expectedStatus: http.StatusBadRequest,
			mockBehavior: func(m *mockTransactionProductUsecase, l *MockLogger) {
				l.On("Error",
					mock.Anything,
					mock.Anything,
				).Return()
			},
		},
		{
			name: "Not Enough Quantity Error",
			inputJSON: `{
                "product_id": "019444a2-e318-79b5-8fe4-b32716306083",
                "product_name": "Product A",
                "quantity": 10,
                "from_warehouse_id": "019444a3-6a3f-7249-b694-f6f071d8eb79",
                "to_warehouse_id": "019444a3-a5dc-7e93-bcc3-fec46dddd299"
            }`,
			expectedStatus: http.StatusInternalServerError,
			mockBehavior: func(m *mockTransactionProductUsecase, l *MockLogger) {
				expectedError := fmt.Errorf("warehouse product quantity is not enough")
				m.On("MoveIn",
					mock.Anything,
					mock.AnythingOfType("*entity.StockMovement"),
				).Return(expectedError)

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
			t.Parallel()

			// initialize mocks
			mockTxUsecase := new(mockTransactionProductUsecase)
			mockStockUsecase := new(mockStockMovementUsecase)
			mockLogger := NewMockLogger(t)

			tt.mockBehavior(mockTxUsecase, mockLogger)

			// setup router
			router := gin.New()
			handler := router.Group("/api/v1")
			newStockMovementRoutes(
				handler,
				mockStockUsecase,
				mockTxUsecase,
				mockLogger,
				func(c *gin.Context) { c.Next() },
			)

			// create request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(
				http.MethodPost,
				"/api/v1/stock-movements/movein",
				bytes.NewBufferString(tt.inputJSON),
			)
			req.Header.Set("Content-Type", "application/json")

			// serve request
			router.ServeHTTP(w, req)

			// assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			mockTxUsecase.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}
