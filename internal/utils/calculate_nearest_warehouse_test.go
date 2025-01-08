package utils

import (
	"testing"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestFindNearestWarehouseByProductID(t *testing.T) {
	t.Parallel()
	warehouse1ID := uuid.New()
	warehouse2ID := uuid.New()
	warehouse3ID := uuid.New()

	tests := []struct {
		name        string
		zipCode     string
		warehouses  []*entity.WarehouseAddressAndProductQty
		expectedZip string
		expectError bool
	}{
		{
			name:    "valid nearest warehouse",
			zipCode: "10000",
			warehouses: []*entity.WarehouseAddressAndProductQty{
				{WarehouseID: warehouse1ID, ZipCode: "12000", ProductName: "Product1", ProductQuantity: 10},
				{WarehouseID: warehouse2ID, ZipCode: "15000", ProductName: "Product1", ProductQuantity: 5},
				{WarehouseID: warehouse3ID, ZipCode: "11000", ProductName: "Product1", ProductQuantity: 15},
			},
			expectedZip: "11000",
			expectError: false,
		},
		{
			name:    "invalid input zipcode",
			zipCode: "abc",
			warehouses: []*entity.WarehouseAddressAndProductQty{
				{WarehouseID: warehouse1ID, ZipCode: "12000", ProductName: "Product1", ProductQuantity: 10},
			},
			expectedZip: "",
			expectError: true,
		},
		{
			name:    "invalid warehouse zipcode",
			zipCode: "10000",
			warehouses: []*entity.WarehouseAddressAndProductQty{
				{WarehouseID: warehouse1ID, ZipCode: "abc", ProductName: "Product1", ProductQuantity: 10},
			},
			expectedZip: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := FindNearestWarehouseByProductID(tt.zipCode, tt.warehouses)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedZip, result)
			}
		})
	}
}

func TestFindNearestWarehouseWithQty(t *testing.T) {
	t.Parallel()
	warehouse1ID := uuid.New()
	warehouse2ID := uuid.New()
	warehouse3ID := uuid.New()

	tests := []struct {
		name        string
		zipCode     string
		warehouses  []*entity.WarehouseAddressAndProductQty
		requestQty  int64
		expected    map[uuid.UUID]int64
		expectError bool
	}{
		{
			name:    "warehouse with sufficient quantity",
			zipCode: "10000",
			warehouses: []*entity.WarehouseAddressAndProductQty{
				{WarehouseID: warehouse1ID, ZipCode: "12345", ProductQuantity: 10},
			},
			requestQty: 5,
			expected: map[uuid.UUID]int64{
				warehouse1ID: 5,
			},
			expectError: false,
		},
		{
			name:    "multiple warehouses needed",
			zipCode: "10000",
			warehouses: []*entity.WarehouseAddressAndProductQty{
				{WarehouseID: warehouse1ID, ZipCode: "11000", ProductQuantity: 5},
				{WarehouseID: warehouse2ID, ZipCode: "12000", ProductQuantity: 5},
				{WarehouseID: warehouse3ID, ZipCode: "13000", ProductQuantity: 5},
			},
			requestQty: 8,
			expected: map[uuid.UUID]int64{
				warehouse1ID: 5,
				warehouse2ID: 3,
			},
			expectError: false,
		},
		{
			name:    "insufficient quantity",
			zipCode: "10000",
			warehouses: []*entity.WarehouseAddressAndProductQty{
				{WarehouseID: warehouse1ID, ZipCode: "11000", ProductQuantity: 5},
				{WarehouseID: warehouse2ID, ZipCode: "12000", ProductQuantity: 3},
			},
			requestQty:  10,
			expected:    nil,
			expectError: true,
		},
		{
			name:    "invalid warehouse zipcode",
			zipCode: "10000",
			warehouses: []*entity.WarehouseAddressAndProductQty{
				{WarehouseID: warehouse1ID, ZipCode: "abc", ProductQuantity: 5},
			},
			requestQty:  5,
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := FindNearestWarehouseWithQty(tt.zipCode, tt.warehouses, tt.requestQty)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
