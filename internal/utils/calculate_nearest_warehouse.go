package utils

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/entity"
)

func FindNearestWarehouse(zipCode string, warehouses []*entity.Warehouse) (string, error) {
	zipCodeNumber, err := strconv.Atoi(zipCode)
	if err != nil {
		return "", fmt.Errorf("invalid zipCodeNumber: %w", err)
	}

	type warehouseEntry struct {
		zipCode  string
		distance int
	}
	entries := make([]warehouseEntry, 0, len(warehouses))
	for _, warehouse := range warehouses {
		warehouseZipCode, err := strconv.Atoi(warehouse.ZipCode)
		if err != nil {
			return "", fmt.Errorf("invalid warehouseZipCode: %w", err)
		}
		entries = append(entries, warehouseEntry{
			zipCode:  warehouse.ZipCode,
			distance: abs(zipCodeNumber - warehouseZipCode),
		})
	}

	// sort by distance value of warehouseDistance
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].distance < entries[j].distance
	})

	return entries[0].zipCode, nil
}

// find nearest warehouse by zipcode
// returned warehouse id with nearest distance (zipcode difference) ascending
func FindNearestWarehouseWithQty(zipcode string, warehouses []*entity.WarehouseAddressAndProductQty, requestQty int64) (map[uuid.UUID]int64, error) {
	zipCodeNumber, err := strconv.Atoi(zipcode)
	if err != nil {
		return nil, fmt.Errorf("invalid zipCodeNumber: %w", err)
	}

	type warehouseEntry struct {
		id       uuid.UUID
		distance int
		quantity int64
	}

	entries := make([]warehouseEntry, 0, len(warehouses))
	for _, warehouse := range warehouses {
		warehouseZipCode, err := strconv.Atoi(warehouse.ZipCode)
		if err != nil {
			return nil, fmt.Errorf("invalid warehouseZipCode: %w", err)
		}
		entries = append(entries, warehouseEntry{
			id:       warehouse.WarehouseID,
			distance: abs(zipCodeNumber - warehouseZipCode),
			quantity: warehouse.ProductQuantity,
		})
	}

	// sort by distance value of warehouseDistance
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].distance < entries[j].distance
	})

	result := make(map[uuid.UUID]int64, len(entries))
	var countLeft int64
	for _, entry := range entries {
		if countLeft >= requestQty {
			break
		}
		result[entry.id] = min(entry.quantity, requestQty-countLeft)
		countLeft += result[entry.id]
	}

	if countLeft < requestQty {
		return nil, fmt.Errorf("not enough product quantity with total quantity %d", requestQty)
	}

	return result, nil
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
