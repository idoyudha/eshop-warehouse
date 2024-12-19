package repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-warehouse/internal/entity"
	rClient "github.com/idoyudha/eshop-warehouse/pkg/redis"
	"github.com/redis/go-redis/v9"
)

const warehouseRankKey = "warehouse:rank"

type WarehouseRankRedisRepo struct {
	*rClient.RedisClient
}

func NewWarehouseRankRedisRepo(client *rClient.RedisClient) *WarehouseRankRedisRepo {
	return &WarehouseRankRedisRepo{
		client,
	}
}

// update distance rankings for warehouse
func (r *WarehouseRankRedisRepo) UpdateWarehouseRanks(ctx context.Context, targetZipcode string, warehouseID string, distances map[string]float64) error {
	pipe := r.Client.Pipeline()

	// add this warehouse's distances to each zipcode's sorted set
	for zipcode, distance := range distances {
		key := fmt.Sprintf("%s:%s", warehouseRankKey, zipcode)
		pipe.ZAdd(ctx, key, redis.Z{
			Score:  distance,
			Member: warehouseID,
		})
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update warehouse ranks: %w", err)
	}

	return nil
}

// gets the ranked list of warehouses for a zipcode
func (r *WarehouseRankRedisRepo) GetNearestWarehouses(ctx context.Context, zipcode string) ([]entity.WarehouseDistance, error) {
	key := fmt.Sprintf("%s:%s", warehouseRankKey, zipcode)
	result, err := r.Client.ZRangeWithScores(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get warehouse ranks: %w", err)
	}

	warehouses := make([]entity.WarehouseDistance, len(result))
	for i, z := range result {
		warehouses[i] = entity.WarehouseDistance{
			WarehouseID: z.Member.(string),
			Distance:    z.Score,
		}
	}

	return warehouses, nil
}

// copies rankings from one warehouse to another (for same zipcode)
func (r *WarehouseRankRedisRepo) CopyRankings(ctx context.Context, sourceWarehouseID, targetWarehouseID uuid.UUID) error {
	pattern := "warehouse:rank:*"
	iter := r.Client.Scan(ctx, 0, pattern, 0).Iterator()

	pipe := r.Client.Pipeline()
	for iter.Next(ctx) {
		key := iter.Val()
		// get score of source warehouse
		score, err := r.Client.ZScore(ctx, key, sourceWarehouseID.String()).Result()
		if err != nil && err != redis.Nil {
			return fmt.Errorf("failed to get source warehouse score: %w", err)
		}

		if err != redis.Nil {
			// add target warehouse with same score
			pipe.ZAdd(ctx, key, redis.Z{
				Score:  score,
				Member: targetWarehouseID.String(),
			})
		}
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to scan redis keys: %w", err)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to copy warehouse rankings: %w", err)
	}

	return nil
}
