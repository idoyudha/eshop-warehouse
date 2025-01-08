package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateWarehouseProductID(t *testing.T) {
	warehouseID, _ := uuid.NewV7()
	productID, _ := uuid.NewV7()
	categoryID, _ := uuid.NewV7()

	tests := []struct {
		name      string
		wp        *WarehouseProduct
		wantPanic bool
	}{
		{
			name: "generate id for empty warehouse product",
			wp:   &WarehouseProduct{},
		},
		{
			name: "generate id for filled warehouse product",
			wp: &WarehouseProduct{
				WarehouseID:        warehouseID,
				ProductID:          productID,
				ProductSKU:         "SKU123",
				ProductName:        "Test Product",
				ProductImageURL:    "http://awss3.com/car.jpg",
				ProductDescription: "Test Description",
				ProductPrice:       14500,
				ProductQuantity:    32,
				ProductCategoryID:  categoryID,
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
				DeletedAt:          time.Time{},
			},
		},
		{
			name:      "should panic for nil warehouse product",
			wp:        nil,
			wantPanic: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantPanic {
				assert.Panics(t, func() {
					_ = tc.wp.GenerateWarehouseProductID()
				})
				return
			}

			var warehouseProduct WarehouseProduct
			if tc.wp != nil {
				warehouseProduct = *tc.wp
			}

			err := tc.wp.GenerateWarehouseProductID()

			assert.NoError(t, err)

			assert.NotEqual(t, uuid.Nil, tc.wp.ID)
			assert.Equal(t, uuid.Version(7), tc.wp.ID.Version())

			if tc.wp != nil {
				assert.Equal(t, warehouseProduct.WarehouseID, tc.wp.WarehouseID)
				assert.Equal(t, warehouseProduct.ProductID, tc.wp.ProductID)
				assert.Equal(t, warehouseProduct.ProductSKU, tc.wp.ProductSKU)
				assert.Equal(t, warehouseProduct.ProductName, tc.wp.ProductName)
				assert.Equal(t, warehouseProduct.ProductImageURL, tc.wp.ProductImageURL)
				assert.Equal(t, warehouseProduct.ProductDescription, tc.wp.ProductDescription)
				assert.Equal(t, warehouseProduct.ProductPrice, tc.wp.ProductPrice)
				assert.Equal(t, warehouseProduct.ProductQuantity, tc.wp.ProductQuantity)
				assert.Equal(t, warehouseProduct.ProductCategoryID, tc.wp.ProductCategoryID)
				assert.Equal(t, warehouseProduct.CreatedAt, tc.wp.CreatedAt)
				assert.Equal(t, warehouseProduct.UpdatedAt, tc.wp.UpdatedAt)
				assert.Equal(t, warehouseProduct.DeletedAt, tc.wp.DeletedAt)
			}
		})
	}
}
