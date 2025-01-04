package v1

import (
	"net/http"

	"github.com/gin-contrib/cors"
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
	handler.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}))

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
