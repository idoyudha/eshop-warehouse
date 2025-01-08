package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateStockMovementID(t *testing.T) {
	productID, _ := uuid.NewV7()
	fromWHID, _ := uuid.NewV7()
	toWHID, _ := uuid.NewV7()
	toUserID, _ := uuid.NewV7()
	tests := []struct {
		name    string
		sm      *StockMovement
		wantErr bool
	}{
		{
			name: "generate id for empty stock movement",
			sm:   &StockMovement{},
		},
		{
			name: "generate id for filled stock movement",
			sm: &StockMovement{
				ProductID:       productID,
				ProductName:     "Test Product",
				Quantity:        10,
				FromWarehouseID: fromWHID,
				ToWarehouseID:   toWHID,
				ToUserID:        toUserID,
				CreatedAt:       time.Now(),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.sm.GenerateStockMovementID()

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, tc.sm.ID)
			assert.Equal(t, uuid.Version(7), tc.sm.ID.Version())
		})
	}
}
