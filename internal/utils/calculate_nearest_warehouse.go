package utils

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/entity"
)

// find nearest warehouse by zipcode
// returned warehouse id with nearest distance (zipcode difference) ascending
func FindNearestWarehouse(zipcode string, warehouses []*entity.WarehouseAddressAndProductQty, requestQty int64) (map[uuid.UUID]int64, error) {
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
		entries = append(entries, warehouseEntry{id: warehouse.WarehouseID, distance: abs(zipCodeNumber - warehouseZipCode)})
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

	return result, nil
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
