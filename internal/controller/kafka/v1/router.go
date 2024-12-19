package v1

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/usecase"
	kafkaConSrv "github.com/idoyudha/eshop-warehouse/pkg/kafka"
	"github.com/idoyudha/eshop-warehouse/pkg/logger"
)

type kafkaConsumerRoutes struct {
	ucp usecase.WarehouseProduct
	l   logger.Interface
}

func KafkaNewRouter(
	ucp usecase.WarehouseProduct,
	l logger.Interface,
	c *kafkaConSrv.ConsumerServer,
) error {
	return nil
}

type KafkaProductUpdatedMessage struct {
	ProductID    uuid.UUID `json:"product_id"`
	ProductName  string    `json:"product_name"`
	ProductPrice float64   `json:"product_price"`
}

func (r *kafkaConsumerRoutes) handleProductUpdated(msg *kafka.Message) error {
	var message KafkaProductUpdatedMessage

	if err := json.Unmarshal(msg.Value, &message); err != nil {
		r.l.Error(err, "http - v1 - kafkaConsumerRoutes - handleProductUpdated")
		return err
	}

	return nil
}
