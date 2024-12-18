package entity

import (
	"time"

	"github.com/google/uuid"
)

type Warehouse struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Street          string    `json:"street"`
	City            string    `json:"city"`
	State           string    `json:"state"`
	ZipCode         string    `json:"zip_code"`
	IsMainWarehouse bool      `json:"is_main_warehouse"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       time.Time `json:"deleted_at"`
}
