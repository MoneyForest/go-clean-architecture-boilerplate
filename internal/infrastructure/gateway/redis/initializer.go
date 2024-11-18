package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Environment string
	Host        string
	Port        string
	Password    string
}

// InitRedis はRedisデータベース接続を初期化します
func InitRedis(ctx context.Context, config RedisConfig) (*redis.Client, error) {
	var options *redis.Options

	switch config.Environment {
	case "local":
		options = &redis.Options{
			Addr:         fmt.Sprintf("%s:%s", config.Host, config.Port),
			Password:     config.Password,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolSize:     10,
			PoolTimeout:  4 * time.Second,
			MaxRetries:   3,
			MinIdleConns: 5,
		}
	default:
		return nil, fmt.Errorf("invalid environment: %s", config.Environment)
	}

	// Redisクライアントを作成
	client := redis.NewClient(options)

	// 接続テスト
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	log.Printf("Successfully connected to Redis at %s:%s\n", config.Host, config.Port)

	return client, nil
}
