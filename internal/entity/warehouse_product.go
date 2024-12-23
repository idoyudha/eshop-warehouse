package entity

import (
	"time"

	"github.com/google/uuid"
)

type WarehouseProduct struct {
	ID                 uuid.UUID `json:"id"`
	WarehouseID        uuid.UUID `json:"warehouse_id"`
	ProductID          uuid.UUID `json:"product_id"`
	ProductSKU         string    `json:"product_sku"`
	ProductName        string    `json:"product_name"`
	ProductImageURL    string    `json:"product_image_url"`
	ProductDescription string    `json:"product_description"`
	ProductPrice       float64   `json:"product_price"`
	ProductQuantity    int64     `json:"product_quantity"`
	ProductCategoryID  uuid.UUID `json:"product_category_id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	DeletedAt          time.Time `json:"deleted_at"`
}

func (wp *WarehouseProduct) GenerateWarehouseProductID() error {
	WarehouseProductID, err := uuid.NewV7()
	if err != nil {
		return err
	}

	wp.ID = WarehouseProductID
	return nil
}

type WarehouseAddressAndProductQty struct {
	WarehouseID     uuid.UUID
	ZipCode         string
	ProductName     string
	ProductQuantity int64
}
