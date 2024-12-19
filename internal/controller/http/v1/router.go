package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idoyudha/eshop-warehouse/config"
	"github.com/idoyudha/eshop-warehouse/internal/usecase"
	"github.com/idoyudha/eshop-warehouse/pkg/logger"
)

func NewRouter(
	handler *gin.Engine,
	ucw usecase.Warehouse,
	ucwp usecase.WarehouseProduct,
	ucsm usecase.StockMovement,
	uct usecase.TransactionProduct,
	l logger.Interface,
	auth config.AuthService,
) {

	handler.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	authMid := cognitoMiddleware(auth)

	h := handler.Group("/v1")
	{
		newWarehouseRoutes(h, ucw, l, authMid)
		newWarehouseProductRoutes(h, ucwp, l, authMid)
		newStockMovementRoutes(h, ucsm, uct, l, authMid)
	}
}
