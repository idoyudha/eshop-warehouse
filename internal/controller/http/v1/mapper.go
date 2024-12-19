package v1

import (
	"time"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/entity"
)

func CreateWarehouseRequestToWarehouseEntity(req CreateWarehouseRequest) (entity.Warehouse, error) {
	warehouseID, err := uuid.NewV7()
	if err != nil {
		return entity.Warehouse{}, err
	}
	return entity.Warehouse{
		ID:              warehouseID,
		Name:            req.Name,
		Street:          req.Street,
		City:            req.City,
		State:           req.State,
		ZipCode:         req.ZipCode,
		IsMainWarehouse: false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil
}

func UpdateWarehouseRequestToWarehouseEntity(req UpdateWarehouseRequest, warehouseID uuid.UUID) (entity.Warehouse, error) {
	return entity.Warehouse{
		ID:        warehouseID,
		Name:      req.Name,
		Street:    req.Street,
		UpdatedAt: time.Now(),
	}, nil
}
