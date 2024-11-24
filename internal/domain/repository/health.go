package repository

import (
	"context"
)

type HealthRepository interface {
	Ping(ctx context.Context) error
}

type HealthCacheRepository interface {
	Ping(ctx context.Context) error
}
