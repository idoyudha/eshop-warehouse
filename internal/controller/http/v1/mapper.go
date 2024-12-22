package v1

import (
	"time"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/entity"
)

func createWarehouseRequestToWarehouseEntity(req createWarehouseRequest) entity.Warehouse {
	return entity.Warehouse{
		Name:            req.Name,
		Street:          req.Street,
		City:            req.City,
		State:           req.State,
		ZipCode:         req.ZipCode,
		IsMainWarehouse: false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func warehouseEntityToCreateWarehouseResponse(warehouse entity.Warehouse) createWarehouseResponse {
	return createWarehouseResponse{
		ID:              warehouse.ID,
		Name:            warehouse.Name,
		Street:          warehouse.Street,
		City:            warehouse.City,
		State:           warehouse.State,
		ZipCode:         warehouse.ZipCode,
		IsMainWarehouse: warehouse.IsMainWarehouse,
	}
}

func warehouseEntitiesToGetAllWarehouseResponse(warehouses []*entity.Warehouse) []getWarehouseResponse {
	var warehouseResponses []getWarehouseResponse
	for _, warehouse := range warehouses {
		warehouseResponses = append(warehouseResponses, warehouseEntityToGetWarehouseResponse(*warehouse))
	}
	return warehouseResponses
}

func updateWarehouseRequestToWarehouseEntity(req updateWarehouseRequest, warehouseID uuid.UUID) entity.Warehouse {
	return entity.Warehouse{
		ID:        warehouseID,
		Name:      req.Name,
		Street:    req.Street,
		UpdatedAt: time.Now(),
	}
}

func warehouseEntityToUpdateWarehouseResponse(warehouse entity.Warehouse) updateWarehouseResponse {
	return updateWarehouseResponse{
		ID:              warehouse.ID,
		Name:            warehouse.Name,
		Street:          warehouse.Street,
		City:            warehouse.City,
		State:           warehouse.State,
		ZipCode:         warehouse.ZipCode,
		IsMainWarehouse: warehouse.IsMainWarehouse,
	}
}

func createStockMovementInRequestToStockMovementEntity(req CreateStockMovementIn) entity.StockMovement {
	return entity.StockMovement{
		ProductID:       req.ProductID,
		ProductName:     req.ProductName,
		Quantity:        req.Quantity,
		FromWarehouseID: req.FromWarehouseID,
		ToWarehouseID:   req.ToWarehouseID,
		CreatedAt:       time.Now(),
	}
}

func warehouseEntityToGetWarehouseResponse(warehouse entity.Warehouse) getWarehouseResponse {
	return getWarehouseResponse{
		ID:              warehouse.ID,
		Name:            warehouse.Name,
		Street:          warehouse.Street,
		City:            warehouse.City,
		State:           warehouse.State,
		ZipCode:         warehouse.ZipCode,
		IsMainWarehouse: warehouse.IsMainWarehouse,
	}
}

func createStockMovementOutRequestToStockMovementEntity(req createStockMovementOut) entity.StockMovement {
	return entity.StockMovement{
		ProductID:       req.ProductID,
		ProductName:     req.ProductName,
		Quantity:        req.Quantity,
		FromWarehouseID: req.FromWarehouseID,
		ToUserID:        req.ToUserID,
		CreatedAt:       time.Now(),
	}
}
