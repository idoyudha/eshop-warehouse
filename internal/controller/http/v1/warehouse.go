package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/usecase"
	"github.com/idoyudha/eshop-warehouse/pkg/logger"
)

type warehouseRoutes struct {
	uc usecase.Warehouse
	l  logger.Interface
}

type CreateWarehouseRequest struct {
	Name    string `json:"name" binding:"required"`
	Street  string `json:"street" binding:"required"`
	City    string `json:"city" binding:"required"`
	State   string `json:"state" binding:"required"`
	ZipCode string `json:"zip_code" binding:"required"`
}

func (r *warehouseRoutes) createWarehouse(ctx *gin.Context) {
	var req CreateWarehouseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - warehouseRoutes - createWarehouse")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	warehouse, err := CreateWarehouseRequestToWarehouseEntity(req)
	if err != nil {
		r.l.Error(err, "http - v1 - warehouseRoutes - createWarehouse")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	err = r.uc.CreateWarehouse(ctx.Request.Context(), &warehouse)
	if err != nil {
		r.l.Error(err, "http - v1 - warehouseRoutes - createWarehouse")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, newCreateSuccess(warehouse))
}

type UpdateWarehouseRequest struct {
	Name   string `json:"name" binding:"required"`
	Street string `json:"street" binding:"required"`
}

func (r *warehouseRoutes) updateWarehouse(ctx *gin.Context) {
	warehouseID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - warehouseRoutes - updateCart")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	var req UpdateWarehouseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - warehouseRoutes - updateWarehouse")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	warehouse, err := UpdateWarehouseRequestToWarehouseEntity(req, warehouseID)
	if err != nil {
		r.l.Error(err, "http - v1 - warehouseRoutes - updateWarehouse")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	err = r.uc.UpdateWarehouse(ctx.Request.Context(), &warehouse)
	if err != nil {
		r.l.Error(err, "http - v1 - warehouseRoutes - updateWarehouse")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, newUpdateSuccess(warehouse))
}

func (r *warehouseRoutes) getWarehouseByID(ctx *gin.Context) {
	warehouseID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - warehouseRoutes - updateCart")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	warehouse, err := r.uc.GetWarehouseByID(ctx.Request.Context(), warehouseID)
	if err != nil {
		r.l.Error(err, "http - v1 - warehouseRoutes - getWarehouseByID")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, newGetSuccess(warehouse))
}

func (r *warehouseRoutes) getAllWarehouses(ctx *gin.Context) {
	warehouses, err := r.uc.GetAllWarehouses(ctx.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - warehouseRoutes - getAllWarehouses")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, newGetSuccess(warehouses))
}
