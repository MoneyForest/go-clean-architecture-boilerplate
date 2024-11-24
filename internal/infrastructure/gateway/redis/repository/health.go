package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type HealthRedisRepository struct {
	client *redis.Client
}

func NewHealthRedisRepository(client *redis.Client) HealthRedisRepository {
	return HealthRedisRepository{client: client}
}

func (c HealthRedisRepository) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}
