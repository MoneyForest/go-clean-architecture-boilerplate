package repository

import (
	"context"
	"time"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
)

type UserRepository interface {
	Save(ctx context.Context, user *model.User) (*model.User, error)
	FindById(ctx context.Context, id uuid.UUID) (*model.User, error)
	FindAll(ctx context.Context, limit, offset int) ([]*model.User, error)
	Remove(ctx context.Context, id uuid.UUID) (*uuid.UUID, error)
}

type UserCacheRepository interface {
	FindById(ctx context.Context, id uuid.UUID) (*model.User, error)
	Store(ctx context.Context, user *model.User, ttl time.Duration) error
	Remove(ctx context.Context, id uuid.UUID) error
}
