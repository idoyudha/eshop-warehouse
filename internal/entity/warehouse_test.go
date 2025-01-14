package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateWarehouseID(t *testing.T) {
	tests := []struct {
		name      string
		warehouse *Warehouse
		wantPanic bool
	}{
		{
			name:      "generate id for empty warehouse",
			warehouse: &Warehouse{},
		},
		{
			name: "generate id for filled warehouse",
			warehouse: &Warehouse{
				Name:            "Test Warehouse",
				Street:          "Bintaro Street",
				City:            "South Tangerang",
				State:           "Banten",
				ZipCode:         "12329",
				IsMainWarehouse: true,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
				DeletedAt:       time.Time{},
			},
		},
		{
			name:      "should panic for nil warehouse",
			warehouse: nil,
			wantPanic: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantPanic {
				assert.Panics(t, func() {
					_ = tc.warehouse.GenerateWarehouseID()
				})
				return
			}

			var wh Warehouse
			if tc.warehouse != nil {
				wh = *tc.warehouse
			}

			err := tc.warehouse.GenerateWarehouseID()

			assert.NoError(t, err)

			assert.NotEqual(t, uuid.Nil, tc.warehouse.ID)
			assert.Equal(t, uuid.Version(7), tc.warehouse.ID.Version())

			if tc.warehouse != nil {
				assert.Equal(t, wh.Name, tc.warehouse.Name)
				assert.Equal(t, wh.Street, tc.warehouse.Street)
				assert.Equal(t, wh.City, tc.warehouse.City)
				assert.Equal(t, wh.State, tc.warehouse.State)
				assert.Equal(t, wh.ZipCode, tc.warehouse.ZipCode)
				assert.Equal(t, wh.IsMainWarehouse, tc.warehouse.IsMainWarehouse)
				assert.Equal(t, wh.CreatedAt, tc.warehouse.CreatedAt)
				assert.Equal(t, wh.UpdatedAt, tc.warehouse.UpdatedAt)
				assert.Equal(t, wh.DeletedAt, tc.warehouse.DeletedAt)
			}
		})
	}
}
