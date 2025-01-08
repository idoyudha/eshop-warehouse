package usecase_test

// import (
// 	"context"
// 	"testing"

// 	"github.com/idoyudha/eshop-warehouse/internal/usecase"
// 	"github.com/idoyudha/eshop-warehouse/pkg/kafka"
// 	gomock "go.uber.org/mock/gomock"
// )

type TestTransferProduct struct {
	name string
	mock func()
	res  error
}

// TODO: mock kafka producer, add to interface

// type MockKafkaProducer struct {
// 	mock *gomock.Controller
// }

// func NewMockKafkaProducer(mock *gomock.Controller) *MockKafkaProducer {
// 	return &MockKafkaProducer{
// 		mock: mock,
// 	}
// }

// func (m *MockKafkaProducer) Produce(topic string, key []byte, value interface{}) error {
// 	return nil
// }

// func transactionProduct(t *testing.T) (
// 	*usecase.TransactionProductUseCase,
// 	*MockWarehouseRankRedisRepo,
// 	*MockTransactionProductPostgresRepo,
// 	*MockWarehouseProductPostgreRepo,
// 	*MockKafkaProducer,
// ) {
// 	t.Helper()

// 	mockCtl := gomock.NewController(t)
// 	defer mockCtl.Finish()

// 	repoRedis := NewMockWarehouseRankRedisRepo(mockCtl)
// 	repoTransactionPostgres := NewMockTransactionProductPostgresRepo(mockCtl)
// 	repoProductPostgres := NewMockWarehouseProductPostgreRepo(mockCtl)
// 	producer := NewMockKafkaProducer(mockCtl)

// 	transactionProduct := usecase.NewTransactionProductUseCase(
// 		repoRedis,
// 		repoTransactionPostgres,
// 		repoProductPostgres,
// 		producer,
// 	)

// 	return transactionProduct, repoRedis, repoTransactionPostgres, repoProductPostgres, producer
// }

// func TestMovementIn(t *testing.T) {
// 	// allow this function run in parallel with other test function
// 	t.Parallel()

// 	tests := []TestTransferProduct{
// 		{
// 			name: "success",
// 			mock: func(repo *MockTransactionProductPostgresRepo) {
// 				repo.EXPECT().
// 					TransferIn(context.Background(), mockStockMovements).
// 					Return(mockStockMovements, nil)
// 			},
// 			res: nil,
// 		},
// 	}

// 	for _, tc := range tests {
// 		tc := tc
// 		t.Run(tc.name, func(t *testing.T) {
// 			// test case will run in parallel
// 			t.Parallel()
// 		})
// 	}
// }
