package entity

import (
	"time"

	"github.com/google/uuid"
)

type StockMovement struct {
	ID              uuid.UUID `json:"id"`
	ProductID       uuid.UUID `json:"product_id"`
	ProductName     string    `json:"product_name"`
	Quantity        int64     `json:"quantity"`
	FromWarehouseID uuid.UUID `json:"from_warehouse_id"`
	ToWarehouseID   uuid.UUID `json:"to_warehouse_id"`
	ToUserID        uuid.UUID `json:"to_user_id"` // for moving out to user (DELIVERED)
	CreatedAt       time.Time `json:"created_at"`
}
