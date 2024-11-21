package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/gateway/redis/dto"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
)

const userKeyPrefix = "user:"

type UserRedisRepository struct {
	client *redis.Client
}

func NewUserRedisRepository(client *redis.Client) UserRedisRepository {
	return UserRedisRepository{client: client}
}

func (c UserRedisRepository) SetWithTTL(ctx context.Context, user *model.User, ttl time.Duration) error {
	entity := dto.ToUserEntity(user)
	jsonData, err := dto.ToJSON(entity)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	key := userKeyPrefix + user.ID.String()
	if err := c.client.Set(ctx, key, jsonData, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set user cache: %w", err)
	}

	return nil
}

func (c UserRedisRepository) Get(ctx context.Context, id uuid.UUID) (*model.User, error) {
	key := userKeyPrefix + id.String()
	jsonData, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("user not found in cache: %s", id)
		}
		return nil, fmt.Errorf("failed to get user from cache: %w", err)
	}

	entity, err := dto.FromJSON(jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return dto.ToUserModel(entity)
}

func (c UserRedisRepository) Delete(ctx context.Context, id uuid.UUID) error {
	key := userKeyPrefix + id.String()
	return c.client.Del(ctx, key).Err()
}
