package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/usecase"
	"github.com/idoyudha/eshop-warehouse/pkg/logger"
)

type stockMovementRoutes struct {
	ucs usecase.StockMovement
	uct usecase.TransactionProduct
	l   logger.Interface
}

func newStockMovementRoutes(
	handler *gin.RouterGroup,
	ucs usecase.StockMovement,
	uct usecase.TransactionProduct,
	l logger.Interface,
	authMid gin.HandlerFunc,
) {
	r := &stockMovementRoutes{
		ucs: ucs,
		uct: uct,
		l:   l,
	}

	h := handler.Group("/stock-movements").Use(authMid)
	{
		h.POST("/movein", r.createStockMovementIn)
		h.POST("/moveout", r.createStockMovementOut)
		h.GET("", r.getAllStockMovements)
		h.GET("/product/:product_id", r.getStockMovementByProductID)
		h.GET("/source/:source_id", r.getStockMovementBySourceID)
		h.GET("/destination/:destination_id", r.getStockMovementByDestinationID)
		// TODO: route for get stock movement destination user id
	}
}

type CreateStockMovementIn struct {
	ProductID       uuid.UUID `json:"product_id"`
	ProductName     string    `json:"product_name"`
	Quantity        int64     `json:"quantity"`
	FromWarehouseID uuid.UUID `json:"from_warehouse_id"`
	ToWarehouseID   uuid.UUID `json:"to_warehouse_id"`
}

func (r *stockMovementRoutes) createStockMovementIn(ctx *gin.Context) {
	var req CreateStockMovementIn
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - stockMovementRoutes - createStockMovementIn")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	stockMovement := createStockMovementInRequestToStockMovementEntity(req)

	err := r.uct.MoveIn(ctx.Request.Context(), &stockMovement)
	if err != nil {
		r.l.Error(err, "http - v1 - stockMovementRoutes - createStockMovementIn")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, newCreateSuccess(stockMovement))
}

type createStockMovementOut struct {
	Items   []itemStockMovementOut `json:"items" binding:"required"`
	ZipCode string                 `json:"zipcode" binding:"required"`
}

type itemStockMovementOut struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int64     `json:"quantity" binding:"required"`
	Price     float64   `json:"price" binding:"required"`
}

func (r *stockMovementRoutes) createStockMovementOut(ctx *gin.Context) {
	var req createStockMovementOut
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - stockMovementRoutes - createStockMovementOut")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	userID, exist := ctx.Get(UserIDKey)
	if !exist {
		r.l.Error("not exist", "http - v1 - stockMovementRoutes - createStockMovementOut")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError("user id not exist"))
		return
	}

	stockMovements := createStockMovementOutRequestToStockMovementEntity(req, userID.(uuid.UUID))
	err := r.uct.MoveOut(ctx.Request.Context(), stockMovements, req.ZipCode)
	if err != nil {
		r.l.Error(err, "http - v1 - stockMovementRoutes - createStockMovementOut")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, newCreateSuccess(stockMovements))
}

func (r *stockMovementRoutes) getAllStockMovements(ctx *gin.Context) {
	stockMovements, err := r.ucs.GetAllStockMovements(ctx.Request.Context())
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

	stockMovements, err := r.ucs.GetStockMovementsByProductID(ctx.Request.Context(), productID)
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

	stockMovements, err := r.ucs.GetStockMovementsBySourceID(ctx.Request.Context(), sourceID)
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

	stockMovements, err := r.ucs.GetStockMovementsByDestinationID(ctx.Request.Context(), destinationID)
	if err != nil {
		r.l.Error(err, "http - v1 - stockMovementRoutes - getStockMovementByDestinationID")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, newGetSuccess(stockMovements))
}
