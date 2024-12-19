package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/usecase"
	"github.com/idoyudha/eshop-warehouse/pkg/logger"
)

type stockMovementRoutes struct {
	uc usecase.StockMovement
	l  logger.Interface
}

func newStockMovementRoutes(handler *gin.RouterGroup, uc usecase.StockMovement, l logger.Interface) {
	r := &stockMovementRoutes{uc: uc, l: l}

	h := handler.Group("/stock-movements")
	{
		h.GET("", r.getAllStockMovements)
		h.GET("/product/:product_id", r.getStockMovementByProductID)
		h.GET("/source/:source_id", r.getStockMovementBySourceID)
		h.GET("/destination/:source_id", r.getStockMovementByDestinationID)
	}
}

func (r *stockMovementRoutes) getAllStockMovements(ctx *gin.Context) {
	stockMovements, err := r.uc.GetAllStockMovements(ctx.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - stockMovementRoutes - getAllStockMovements")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, newGetSuccess(stockMovements))
}

func (r *stockMovementRoutes) getStockMovementByProductID(ctx *gin.Context) {
	productID, err := uuid.Parse(ctx.Param("product_id"))
	if err != nil {
		r.l.Error(err, "http - v1 - stockMovementRoutes - getStockMovementByProductID")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	stockMovements, err := r.uc.GetStockMovementsByProductID(ctx.Request.Context(), productID)
	if err != nil {
		r.l.Error(err, "http - v1 - stockMovementRoutes - getStockMovementByProductID")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, newGetSuccess(stockMovements))
}

func (r *stockMovementRoutes) getStockMovementBySourceID(ctx *gin.Context) {
	sourceID, err := uuid.Parse(ctx.Param("source_id"))
	if err != nil {
		r.l.Error(err, "http - v1 - stockMovementRoutes - getStockMovementBySourceID")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	stockMovements, err := r.uc.GetStockMovementsBySourceID(ctx.Request.Context(), sourceID)
	if err != nil {
		r.l.Error(err, "http - v1 - stockMovementRoutes - getStockMovementBySourceID")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, newGetSuccess(stockMovements))
}

func (r *stockMovementRoutes) getStockMovementByDestinationID(ctx *gin.Context) {
	destinationID, err := uuid.Parse(ctx.Param("destination_id"))
	if err != nil {
		r.l.Error(err, "http - v1 - stockMovementRoutes - getStockMovementByDestinationID")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	stockMovements, err := r.uc.GetStockMovementsByDestinationID(ctx.Request.Context(), destinationID)
	if err != nil {
		r.l.Error(err, "http - v1 - stockMovementRoutes - getStockMovementByDestinationID")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, newGetSuccess(stockMovements))
}
