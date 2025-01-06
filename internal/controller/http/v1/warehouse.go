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

func newWarehouseRoutes(handler *gin.RouterGroup, uc usecase.Warehouse, l logger.Interface, authMid gin.HandlerFunc) {
	r := &warehouseRoutes{uc: uc, l: l}

	h := handler.Group("/warehouse").Use(authMid)
	{
		h.POST("", r.createWarehouse)
		h.GET("", r.getAllWarehouses)
		h.GET("/:id", r.getWarehouseByID)
		h.PATCH("/:id", r.updateWarehouse)
		h.POST("/nearest", r.getNearestWarehouse)
	}
}

type createWarehouseRequest struct {
	Name    string `json:"name" binding:"required"`
	Street  string `json:"street" binding:"required"`
	City    string `json:"city" binding:"required"`
	State   string `json:"state" binding:"required"`
	ZipCode string `json:"zip_code" binding:"required"`
}

type createWarehouseResponse struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Street          string    `json:"street"`
	City            string    `json:"city"`
	State           string    `json:"state"`
	ZipCode         string    `json:"zip_code"`
	IsMainWarehouse bool      `json:"is_main_warehouse"`
}

func (r *warehouseRoutes) createWarehouse(ctx *gin.Context) {
	var req createWarehouseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - warehouseRoutes - createWarehouse")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	warehouse := createWarehouseRequestToWarehouseEntity(req)

	err := r.uc.CreateWarehouse(ctx.Request.Context(), &warehouse)
	if err != nil {
		r.l.Error(err, "http - v1 - warehouseRoutes - createWarehouse")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	warehouseResponse := warehouseEntityToCreateWarehouseResponse(warehouse)

	ctx.JSON(http.StatusCreated, newCreateSuccess(warehouseResponse))
}

type updateWarehouseRequest struct {
	Name   string `json:"name" binding:"required"`
	Street string `json:"street" binding:"required"`
}

type updateWarehouseResponse struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Street          string    `json:"street"`
	City            string    `json:"city"`
	State           string    `json:"state"`
	ZipCode         string    `json:"zip_code"`
	IsMainWarehouse bool      `json:"is_main_warehouse"`
}

func (r *warehouseRoutes) updateWarehouse(ctx *gin.Context) {
	warehouseID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - warehouseRoutes - updateWarehouse")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	var req updateWarehouseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - warehouseRoutes - updateWarehouse")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	warehouse := updateWarehouseRequestToWarehouseEntity(req, warehouseID)

	err = r.uc.UpdateWarehouse(ctx.Request.Context(), &warehouse)
	if err != nil {
		r.l.Error(err, "http - v1 - warehouseRoutes - updateWarehouse")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	warehouseResponse := warehouseEntityToUpdateWarehouseResponse(warehouse)

	ctx.JSON(http.StatusOK, newUpdateSuccess(warehouseResponse))
}

type getWarehouseResponse struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Street          string    `json:"street"`
	City            string    `json:"city"`
	State           string    `json:"state"`
	ZipCode         string    `json:"zip_code"`
	IsMainWarehouse bool      `json:"is_main_warehouse"`
}

func (r *warehouseRoutes) getWarehouseByID(ctx *gin.Context) {
	warehouseID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - warehouseRoutes - getWarehouseByID")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	warehouse, err := r.uc.GetWarehouseByID(ctx.Request.Context(), warehouseID)
	if err != nil {
		r.l.Error(err, "http - v1 - warehouseRoutes - getWarehouseByID")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	warehouseResponse := warehouseEntityToGetWarehouseResponse(*warehouse)

	ctx.JSON(http.StatusOK, newGetSuccess(warehouseResponse))
}

func (r *warehouseRoutes) getAllWarehouses(ctx *gin.Context) {
	warehouses, err := r.uc.GetAllWarehouses(ctx.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - warehouseRoutes - getAllWarehouses")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	warehousesResponse := warehouseEntitiesToGetAllWarehouseResponse(warehouses)

	ctx.JSON(http.StatusOK, newGetSuccess(warehousesResponse))
}

type getNearestWarehouseRequest struct {
	ZipCodes []string `json:"zip_codes" binding:"required"`
}

type getNearestWarehouseResponse struct {
	Warehouses map[string]string `json:"warehouses"`
}

func (r *warehouseRoutes) getNearestWarehouse(ctx *gin.Context) {
	var req getNearestWarehouseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - warehouseRoutes - createWarehouse")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	nearestWarehouse, err := r.uc.GetNearestWarehouse(ctx.Request.Context(), req.ZipCodes)
	if err != nil {
		r.l.Error(err, "http - v1 - warehouseRoutes - createWarehouse")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	var res getNearestWarehouseResponse
	res.Warehouses = make(map[string]string)
	for zipCodeFrom, zipCodeTo := range nearestWarehouse {
		res.Warehouses[zipCodeFrom] = zipCodeTo
	}

	ctx.JSON(http.StatusOK, newGetSuccess(res))
}
