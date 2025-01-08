package v1

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestCreateWarehouseRequestToWarehouseEntity(t *testing.T) {
	req := createWarehouseRequest{
		Name:    "Test Warehouse",
		Street:  "123 Test St",
		City:    "Test City",
		State:   "Test State",
		ZipCode: "12345",
	}

	result := createWarehouseRequestToWarehouseEntity(req)

	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.Street, result.Street)
	assert.Equal(t, req.City, result.City)
	assert.Equal(t, req.State, result.State)
	assert.Equal(t, req.ZipCode, result.ZipCode)
	assert.False(t, result.IsMainWarehouse)
	assert.WithinDuration(t, time.Now(), result.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), result.UpdatedAt, time.Second)
}

func TestWarehouseEntityToCreateWarehouseResponse(t *testing.T) {
	id := uuid.New()
	warehouse := entity.Warehouse{
		ID:              id,
		Name:            "Test Warehouse",
		Street:          "123 Test St",
		City:            "Test City",
		State:           "Test State",
		ZipCode:         "12345",
		IsMainWarehouse: true,
	}

	result := warehouseEntityToCreateWarehouseResponse(warehouse)

	assert.Equal(t, warehouse.ID, result.ID)
	assert.Equal(t, warehouse.Name, result.Name)
	assert.Equal(t, warehouse.Street, result.Street)
	assert.Equal(t, warehouse.City, result.City)
	assert.Equal(t, warehouse.State, result.State)
	assert.Equal(t, warehouse.ZipCode, result.ZipCode)
	assert.Equal(t, warehouse.IsMainWarehouse, result.IsMainWarehouse)
}

func TestWarehouseEntitiesToGetAllWarehouseResponse(t *testing.T) {
	warehouses := []*entity.Warehouse{
		{
			ID:              uuid.New(),
			Name:            "Warehouse 1",
			Street:          "Street 1",
			City:            "City 1",
			State:           "State 1",
			ZipCode:         "11111",
			IsMainWarehouse: true,
		},
		{
			ID:              uuid.New(),
			Name:            "Warehouse 2",
			Street:          "Street 2",
			City:            "City 2",
			State:           "State 2",
			ZipCode:         "22222",
			IsMainWarehouse: false,
		},
	}

	result := warehouseEntitiesToGetAllWarehouseResponse(warehouses)

	assert.Len(t, result, len(warehouses))
	for i, warehouse := range warehouses {
		assert.Equal(t, warehouse.ID, result[i].ID)
		assert.Equal(t, warehouse.Name, result[i].Name)
		assert.Equal(t, warehouse.Street, result[i].Street)
		assert.Equal(t, warehouse.City, result[i].City)
		assert.Equal(t, warehouse.State, result[i].State)
		assert.Equal(t, warehouse.ZipCode, result[i].ZipCode)
		assert.Equal(t, warehouse.IsMainWarehouse, result[i].IsMainWarehouse)
	}
}

func TestUpdateWarehouseRequestToWarehouseEntity(t *testing.T) {
	id := uuid.New()
	req := updateWarehouseRequest{
		Name:   "Updated Warehouse",
		Street: "Updated Street",
	}

	result := updateWarehouseRequestToWarehouseEntity(req, id)

	assert.Equal(t, id, result.ID)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.Street, result.Street)
	assert.WithinDuration(t, time.Now(), result.UpdatedAt, time.Second)
}

func TestCreateStockMovementInRequestToStockMovementEntity(t *testing.T) {
	req := CreateStockMovementIn{
		ProductID:       uuid.New(),
		ProductName:     "Test Product",
		Quantity:        10,
		FromWarehouseID: uuid.New(),
		ToWarehouseID:   uuid.New(),
	}

	result := createStockMovementInRequestToStockMovementEntity(req)

	assert.Equal(t, req.ProductID, result.ProductID)
	assert.Equal(t, req.ProductName, result.ProductName)
	assert.Equal(t, req.Quantity, result.Quantity)
	assert.Equal(t, req.FromWarehouseID, result.FromWarehouseID)
	assert.Equal(t, req.ToWarehouseID, result.ToWarehouseID)
	assert.WithinDuration(t, time.Now(), result.CreatedAt, time.Second)
}

func TestCreateStockMovementOutRequestToStockMovementEntity(t *testing.T) {
	userID := uuid.New()
	req := createStockMovementOut{
		Items: []ItemStockMovementOut{
			{
				ProductID: uuid.New(),
				Quantity:  5,
			},
			{
				ProductID: uuid.New(),
				Quantity:  10,
			},
		},
	}

	result := createStockMovementOutRequestToStockMovementEntity(req, userID)

	assert.Len(t, result, len(req.Items))
	for i, movement := range result {
		assert.Equal(t, req.Items[i].ProductID, movement.ProductID)
		assert.Equal(t, req.Items[i].Quantity, movement.Quantity)
		assert.Equal(t, userID, movement.ToUserID)
		assert.WithinDuration(t, time.Now(), movement.CreatedAt, time.Second)
	}
}

func TestWarehouseEntityToUpdateWarehouseResponse(t *testing.T) {
	warehouse := entity.Warehouse{
		ID:              uuid.New(),
		Name:            "Test Warehouse",
		Street:          "Test Street",
		City:            "Test City",
		State:           "Test State",
		ZipCode:         "12345",
		IsMainWarehouse: true,
	}

	result := warehouseEntityToUpdateWarehouseResponse(warehouse)

	assert.Equal(t, warehouse.ID, result.ID)
	assert.Equal(t, warehouse.Name, result.Name)
	assert.Equal(t, warehouse.Street, result.Street)
	assert.Equal(t, warehouse.City, result.City)
	assert.Equal(t, warehouse.State, result.State)
	assert.Equal(t, warehouse.ZipCode, result.ZipCode)
	assert.Equal(t, warehouse.IsMainWarehouse, result.IsMainWarehouse)
}
