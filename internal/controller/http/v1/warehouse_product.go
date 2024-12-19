package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/usecase"
	"github.com/idoyudha/eshop-warehouse/pkg/logger"
)

type warehouseProductRoutes struct {
	uc usecase.WarehouseProduct
	l  logger.Interface
}

func newWarehouseProductRoutes(
	handler *gin.RouterGroup,
	uc usecase.WarehouseProduct,
	l logger.Interface,
	authMid gin.HandlerFunc,
) {
	r := &warehouseProductRoutes{uc: uc, l: l}

	h := handler.Group("/warehouse-product").Use(authMid)
	{
		h.GET("", r.getAllWarehouseProducts)
		h.GET("/product/:product_id", r.getWarehouseProductByProductID)
		h.GET("/warehouse/:warehouse_id", r.getWarehouseProductByWarehouseID)
		h.GET("/product/:product_id/warehouse/:warehouse_id", r.getWarehouseProductByProductIDAndWarehouseID)
	}
}

func (r *warehouseProductRoutes) getAllWarehouseProducts(ctx *gin.Context) {
	products, err := r.uc.GetAllWarehouseProducts(ctx.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - warehouseProductRoutes - getAllWarehouseProducts")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, newGetSuccess(products))
}

func (r *warehouseProductRoutes) getWarehouseProductByProductID(ctx *gin.Context) {
	productID, err := uuid.Parse(ctx.Param("product_id"))
	if err != nil {
		r.l.Error(err, "http - v1 - warehouseProductRoutes - getWarehouseProductByProductID")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	products, err := r.uc.GetWarehouseProductByProductID(ctx.Request.Context(), productID)
	if err != nil {
		r.l.Error(err, "http - v1 - warehouseProductRoutes - getWarehouseProductByProductID")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, newGetSuccess(products))
}

func (r *warehouseProductRoutes) getWarehouseProductByWarehouseID(ctx *gin.Context) {
	warehouseID, err := uuid.Parse(ctx.Param("warehouse_id"))
	if err != nil {
		r.l.Error(err, "http - v1 - warehouseProductRoutes - getWarehouseProductByWarehouseID")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	products, err := r.uc.GetWarehouseProductByWarehouseID(ctx.Request.Context(), warehouseID)
	if err != nil {
		r.l.Error(err, "http - v1 - warehouseProductRoutes - getWarehouseProductByWarehouseID")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, newGetSuccess(products))
}

func (r *warehouseProductRoutes) getWarehouseProductByProductIDAndWarehouseID(ctx *gin.Context) {
	productID, err := uuid.Parse(ctx.Param("product_id"))
	if err != nil {
		r.l.Error(err, "http - v1 - warehouseProductRoutes - getWarehouseProductByProductIDAndWarehouseID")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	warehouseID, err := uuid.Parse(ctx.Param("warehouse_id"))
	if err != nil {
		r.l.Error(err, "http - v1 - warehouseProductRoutes - getWarehouseProductByProductIDAndWarehouseID")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	products, err := r.uc.GetWarehouseProductByProductIDAndWarehouseID(ctx.Request.Context(), productID, warehouseID)
	if err != nil {
		r.l.Error(err, "http - v1 - warehouseProductRoutes - getWarehouseProductByProductIDAndWarehouseID")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, newGetSuccess(products))
}
