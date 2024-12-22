package v1

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/entity"
	"github.com/idoyudha/eshop-warehouse/internal/usecase"
	kafkaConSrv "github.com/idoyudha/eshop-warehouse/pkg/kafka"
	"github.com/idoyudha/eshop-warehouse/pkg/logger"
)

type kafkaConsumerRoutes struct {
	ucw usecase.Warehouse
	ucp usecase.WarehouseProduct
	l   logger.Interface
}

func KafkaNewRouter(
	ucw usecase.Warehouse,
	ucp usecase.WarehouseProduct,
	l logger.Interface,
	c *kafkaConSrv.ConsumerServer,
) error {
	routes := &kafkaConsumerRoutes{
		ucw: ucw,
		ucp: ucp,
		l:   l,
	}

	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	run := true
	for run {
		select {
		case sig := <-sigchan:
			log.Printf("Caught signal %v: terminating\n", sig)
			run = false
			return nil
		default:
			ev, err := c.Consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// log.Println("CONSUME CART SERVICE!!")
				// Errors are informational and automatically handled by the consumer
				if err.(kafka.Error).Code() == kafka.ErrTimedOut {
					continue
				}
				l.Error("Error reading message: ", err)
				continue
			}

			switch *ev.TopicPartition.Topic {
			case kafkaConSrv.ProductCreatedTopic:
				if err := routes.handleProductCreated(ev); err != nil {
					l.Error("Failed to handle product creation: %w", err)
				}
			case kafkaConSrv.ProductUpdatedTopic:
				if err := routes.handleProductUpdated(ev); err != nil {
					l.Error("Failed to handle product update: %w", err)
				}
			default:
				l.Info("Unknown topic: %s", *ev.TopicPartition.Topic)
			}

			log.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
				*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
		}
	}

	return nil
}

type kafkaProductCreatedMessage struct {
	ID          uuid.UUID `json:"id"`
	SKU         string    `json:"sku"`
	Name        string    `json:"name"`
	ImageURL    string    `json:"image_url"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	CategoryID  uuid.UUID `json:"category_id"`
}

func (r *kafkaConsumerRoutes) handleProductCreated(msg *kafka.Message) error {
	var message kafkaProductCreatedMessage

	if err := json.Unmarshal(msg.Value, &message); err != nil {
		r.l.Error(err, "http - v1 - kafkaConsumerRoutes - handleProductCreated")
		return err
	}

	warehouseMainID, err := r.ucw.GetMainIDWarehouse(context.Background())
	if err != nil {
		r.l.Error(err, "http - v1 - kafkaConsumerRoutes - handleProductCreated")
		return err
	}

	product := &entity.WarehouseProduct{
		WarehouseID:        warehouseMainID,
		ProductID:          message.ID,
		ProductSKU:         message.SKU,
		ProductName:        message.Name,
		ProductImageURL:    message.ImageURL,
		ProductDescription: message.Description,
		ProductPrice:       message.Price,
		ProductQuantity:    int64(message.Quantity),
		ProductCategoryID:  message.CategoryID,
	}

	if err := r.ucp.CreateWarehouseProduct(context.Background(), product); err != nil {
		r.l.Error(err, "http - v1 - kafkaConsumerRoutes - handleProductCreated")
		return err
	}

	r.l.Info("Product created", "http - v1 - kafkaConsumerRoutes - handleProductCreated")

	return nil
}

type kafkaProductUpdatedMessage struct {
	ProductID          uuid.UUID `json:"product_id"`
	ProductName        string    `json:"product_name"`
	ProductImageURL    string    `json:"product_image_url"`
	ProductDescription string    `json:"product_description"`
	ProductPrice       float64   `json:"product_price"`
	ProductCategoryID  uuid.UUID `json:"product_category_id"`
}

func (r *kafkaConsumerRoutes) handleProductUpdated(msg *kafka.Message) error {
	var message kafkaProductUpdatedMessage

	if err := json.Unmarshal(msg.Value, &message); err != nil {
		r.l.Error(err, "http - v1 - kafkaConsumerRoutes - handleProductUpdated")
		return err
	}

	product := &entity.WarehouseProduct{
		ProductID:          message.ProductID,
		ProductName:        message.ProductName,
		ProductImageURL:    message.ProductImageURL,
		ProductDescription: message.ProductDescription,
		ProductPrice:       message.ProductPrice,
		ProductCategoryID:  message.ProductCategoryID,
	}

	if err := r.ucp.UpdateWarehouseProduct(context.Background(), product); err != nil {
		r.l.Error(err, "http - v1 - kafkaConsumerRoutes - handleProductUpdated")
		return err
	}

	r.l.Info("Product updated", "http - v1 - kafkaConsumerRoutes - handleProductUpdated")

	return nil
}
