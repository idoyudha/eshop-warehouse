package v1

import (
	"bytes"
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

type mockWarehouseUsecase struct {
	mock.Mock
}

func (m *mockWarehouseUsecase) CreateWarehouse(ctx context.Context, warehouse *entity.Warehouse) error {
	args := m.Called(ctx, warehouse)
	return args.Error(0)
}

func (m *mockWarehouseUsecase) UpdateWarehouse(ctx context.Context, warehouse *entity.Warehouse) error {
	args := m.Called(ctx, warehouse)
	return args.Error(0)
}

func (m *mockWarehouseUsecase) GetWarehouseByID(ctx context.Context, id uuid.UUID) (*entity.Warehouse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Warehouse), args.Error(1)
}

func (m *mockWarehouseUsecase) GetAllWarehouses(ctx context.Context) ([]*entity.Warehouse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Warehouse), args.Error(1)
}

func (m *mockWarehouseUsecase) GetMainIDWarehouse(ctx context.Context) (uuid.UUID, error) {
	args := m.Called(ctx)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *mockWarehouseUsecase) GetNearestWarehouse(ctx context.Context, zipCodes []string) (map[string]string, error) {
	args := m.Called(ctx, zipCodes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]string), args.Error(1)
}

func TestCreateWarehouse(t *testing.T) {
	// t.Parallell()

	tests := []struct {
		name         string
		inputJSON    string
		expectedCode int
		setupMock    func(*mockWarehouseUsecase, *MockLogger)
	}{
		{
			name: "success - Warehouse A",
			inputJSON: `{
                "name": "Warehouse A",
                "street": "Carissa Bintaro Street",
                "city": "Tangerang Selatan",
                "state": "Banten",
                "zip_code": "12345"
            }`,
			expectedCode: http.StatusCreated,
			setupMock: func(m *mockWarehouseUsecase, l *MockLogger) {
				m.On("CreateWarehouse",
					mock.Anything,
					mock.MatchedBy(func(w *entity.Warehouse) bool {
						return w.Name == "Warehouse A" &&
							w.Street == "Carissa Bintaro Street" &&
							w.City == "Tangerang Selatan" &&
							w.State == "Banten" &&
							w.ZipCode == "12345"
					}),
				).Return(nil)
			},
		},
		{
			name: "invalid request - missing required fields",
			inputJSON: `{
                "name": "Warehouse A"
            }`,
			expectedCode: http.StatusBadRequest,
			setupMock: func(m *mockWarehouseUsecase, l *MockLogger) {
				l.On("Error", mock.Anything, mock.Anything).Return()
			},
		},
		{
			name: "internal server error - database error",
			inputJSON: `{
                "name": "Warehouse A",
                "street": "Carissa Bintaro Street",
                "city": "Tangerang Selatan",
                "state": "Banten",
                "zip_code": "12345"
            }`,
			expectedCode: http.StatusInternalServerError,
			setupMock: func(m *mockWarehouseUsecase, l *MockLogger) {
				m.On("CreateWarehouse", mock.Anything, mock.Anything).Return(fmt.Errorf("database error"))
				l.On("Error", mock.Anything, mock.Anything).Return()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallell()
			mockUC := new(mockWarehouseUsecase)
			mockLogger := NewMockLogger(t)

			tt.setupMock(mockUC, mockLogger)

			router := gin.New()
			handler := router.Group("/api/v1")
			newWarehouseRoutes(
				handler,
				mockUC,
				mockLogger,
				func(c *gin.Context) { c.Next() },
			)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(
				http.MethodPost,
				"/api/v1/warehouse",
				bytes.NewBufferString(tt.inputJSON),
			)
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			if t.Failed() {
				t.Logf("Response Status: %d", w.Code)
				t.Logf("Response Body: %s", w.Body.String())
			}

			assert.Equal(t, tt.expectedCode, w.Code)
			mockUC.AssertExpectations(t)
			mockLogger.AssertExpectations(t)

			if tt.expectedCode == http.StatusCreated {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "data")
			}
		})
	}
}

func TestGetWarehouseByID(t *testing.T) {
	// t.Parallell()

	mockWarehouse := &entity.Warehouse{
		ID:              uuid.New(),
		Name:            "Test Warehouse",
		Street:          "Carissa Bintaro Street",
		City:            "Tangerang Selatan",
		State:           "Banten",
		ZipCode:         "12345",
		IsMainWarehouse: false,
	}

	tests := []struct {
		name         string
		warehouseID  string
		expectedCode int
		setupMock    func(*mockWarehouseUsecase, *MockLogger)
	}{
		{
			name:         "success",
			warehouseID:  mockWarehouse.ID.String(),
			expectedCode: http.StatusOK,
			setupMock: func(m *mockWarehouseUsecase, l *MockLogger) {
				m.On("GetWarehouseByID",
					mock.Anything,
					mockWarehouse.ID,
				).Return(mockWarehouse, nil)
			},
		},
		{
			name:         "invalid uuid format",
			warehouseID:  "invalid-uuid",
			expectedCode: http.StatusBadRequest,
			setupMock: func(m *mockWarehouseUsecase, l *MockLogger) {
				l.On("Error",
					mock.Anything,
					mock.Anything,
				).Return()
			},
		},
		{
			name:         "warehouse not found",
			warehouseID:  uuid.New().String(),
			expectedCode: http.StatusInternalServerError,
			setupMock: func(m *mockWarehouseUsecase, l *MockLogger) {
				m.On("GetWarehouseByID",
					mock.Anything,
					mock.Anything,
				).Return(nil, fmt.Errorf("warehouse not found"))
				l.On("Error",
					mock.Anything,
					mock.Anything,
				).Return()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallell()
			mockUC := new(mockWarehouseUsecase)
			mockLogger := NewMockLogger(t)

			tt.setupMock(mockUC, mockLogger)

			router := gin.New()
			handler := router.Group("/api/v1")
			newWarehouseRoutes(
				handler,
				mockUC,
				mockLogger,
				func(c *gin.Context) { c.Next() },
			)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/api/v1/warehouse/%s", tt.warehouseID),
				nil,
			)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			mockUC.AssertExpectations(t)
			mockLogger.AssertExpectations(t)

			if tt.expectedCode == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				data, ok := response["data"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, mockWarehouse.Name, data["name"])
				assert.Equal(t, mockWarehouse.Street, data["street"])
				assert.Equal(t, mockWarehouse.City, data["city"])
				assert.Equal(t, mockWarehouse.State, data["state"])
				assert.Equal(t, mockWarehouse.ZipCode, data["zip_code"])
			}
		})
	}
}

func TestGetAllWarehouses(t *testing.T) {
	// t.Parallell()

	mockWarehouses := []*entity.Warehouse{
		{
			ID:              uuid.New(),
			Name:            "Warehouse 1",
			Street:          "123 First St",
			City:            "First City",
			State:           "First State",
			ZipCode:         "12345",
			IsMainWarehouse: true,
		},
		{
			ID:              uuid.New(),
			Name:            "Warehouse 2",
			Street:          "456 Second St",
			City:            "Second City",
			State:           "Second State",
			ZipCode:         "67890",
			IsMainWarehouse: false,
		},
	}

	tests := []struct {
		name          string
		expectedCode  int
		setupMock     func(*mockWarehouseUsecase, *MockLogger)
		checkResponse func(*testing.T, []byte)
	}{
		{
			name:         "success",
			expectedCode: http.StatusOK,
			setupMock: func(m *mockWarehouseUsecase, l *MockLogger) {
				m.On("GetAllWarehouses", mock.Anything).Return(mockWarehouses, nil)
			},
			checkResponse: func(t *testing.T, body []byte) {
				var response struct {
					Data []getWarehouseResponse `json:"data"`
				}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Len(t, response.Data, 2)
				for i := range response.Data {
					assert.Equal(t, mockWarehouses[i].Name, response.Data[i].Name)
					assert.Equal(t, mockWarehouses[i].Street, response.Data[i].Street)
					assert.Equal(t, mockWarehouses[i].City, response.Data[i].City)
					assert.Equal(t, mockWarehouses[i].State, response.Data[i].State)
					assert.Equal(t, mockWarehouses[i].ZipCode, response.Data[i].ZipCode)
				}
			},
		},
		{
			name:         "empty warehouse list",
			expectedCode: http.StatusOK,
			setupMock: func(m *mockWarehouseUsecase, l *MockLogger) {
				m.On("GetAllWarehouses", mock.Anything).Return([]*entity.Warehouse{}, nil)
			},
			checkResponse: func(t *testing.T, body []byte) {
				var response struct {
					Data []getWarehouseResponse `json:"data"`
				}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Empty(t, response.Data)
			},
		},
		{
			name:         "database error",
			expectedCode: http.StatusInternalServerError,
			setupMock: func(m *mockWarehouseUsecase, l *MockLogger) {
				m.On("GetAllWarehouses", mock.Anything).Return(nil, fmt.Errorf("database error"))
				l.On("Error", mock.Anything, mock.Anything).Return()
			},
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallell()
			mockUC := new(mockWarehouseUsecase)
			mockLogger := NewMockLogger(t)

			tt.setupMock(mockUC, mockLogger)

			router := gin.New()
			handler := router.Group("/api/v1")
			newWarehouseRoutes(
				handler,
				mockUC,
				mockLogger,
				func(c *gin.Context) { c.Next() },
			)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(
				http.MethodGet,
				"/api/v1/warehouse",
				nil,
			)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.Bytes())
			}

			mockUC.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}
