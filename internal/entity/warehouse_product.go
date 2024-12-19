package entity

import (
	"time"

	"github.com/google/uuid"
)

type WarehouseProduct struct {
	ID              uuid.UUID `json:"id"`
	WarehouseID     uuid.UUID `json:"warehouse_id"`
	ProductID       uuid.UUID `json:"product_id"`
	ProductName     string    `json:"product_name"`
	ProductQuantity int64     `json:"product_quantity"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       time.Time `json:"deleted_at"`
}
