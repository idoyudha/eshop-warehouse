package entity

import (
	"time"

	"github.com/google/uuid"
)

type Warehouse struct {
	ID              uuid.UUID
	Name            string
	Street          string
	City            string
	State           string
	ZipCode         string
	IsMainWarehouse bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       time.Time
}

func (w *Warehouse) GenerateWarehouseID() error {
	warehouseID, err := uuid.NewV7()
	if err != nil {
		return err
	}

	w.ID = warehouseID
	return nil
}
