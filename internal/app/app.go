package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/idoyudha/eshop-warehouse/config"
	v1Http "github.com/idoyudha/eshop-warehouse/internal/controller/http/v1"
	v1Kafka "github.com/idoyudha/eshop-warehouse/internal/controller/kafka/v1"
	"github.com/idoyudha/eshop-warehouse/internal/usecase"
	"github.com/idoyudha/eshop-warehouse/internal/usecase/repo"
	"github.com/idoyudha/eshop-warehouse/pkg/httpserver"
	"github.com/idoyudha/eshop-warehouse/pkg/kafka"
	"github.com/idoyudha/eshop-warehouse/pkg/logger"
	"github.com/idoyudha/eshop-warehouse/pkg/postgresql"
	"github.com/idoyudha/eshop-warehouse/pkg/redis"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	kafkaConsumer, err := kafka.NewKafkaConsumer(cfg.Kafka.Broker)
	if err != nil {
		l.Fatal("app - Run - kafka.NewKafkaConsumer: ", err)
	}
	defer kafkaConsumer.Close()

	postgreSQL, err := postgresql.NewPostgres(cfg.PostgreSQL)
	if err != nil {
		l.Fatal("app - Run - postgresql.NewPostgres: ", err)
	}

	redisClient, err := redis.NewRedis(cfg.Redis)
	if err != nil {
		l.Fatal("app - Run - redis.NewRedis: ", err)
	}

	warehouseUseCase := usecase.NewWarehouseUseCase(
		repo.NewWarehouseRankRedisRepo(redisClient),
		repo.NewWarehousePostgreRepo(postgreSQL),
	)

	warehouseProductUseCase := usecase.NewWarehouseProductUseCase(
		repo.NewWarehouseProductPostgreRepo(postgreSQL),
	)

	stockMovementUseCase := usecase.NewStockMovementUseCase(
		repo.NewStockMovementPostgreRepo(postgreSQL),
	)

	transactionProductUseCase := usecase.NewTransactionProductUseCase(
		repo.NewWarehouseRankRedisRepo(redisClient),
		repo.NewTransactionProductPostgreRepo(postgreSQL),
		repo.NewWarehouseProductPostgreRepo(postgreSQL),
	)

	// HTTP Server
	handler := gin.Default()
	v1Http.NewRouter(handler, warehouseUseCase, warehouseProductUseCase, stockMovementUseCase, transactionProductUseCase, l, cfg.AuthService)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Kafka Consumer
	kafkaErrChan := make(chan error, 1)
	go func() {
		if err := v1Kafka.KafkaNewRouter(warehouseUseCase, warehouseProductUseCase, l, kafkaConsumer); err != nil {
			kafkaErrChan <- err
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error("app - Run - httpServer.Notify: ", err)
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Info("app - Run - httpServer.Shutdown: %s", err)
	}
}
