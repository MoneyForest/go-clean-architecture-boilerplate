package repository

import (
	"context"
	"time"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Get(ctx context.Context, id uuid.UUID) (*model.User, error)
	List(ctx context.Context, limit, offset int) ([]*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id uuid.UUID) (*uuid.UUID, error)
}

type UserCacheRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*model.User, error)
	SetWithTTL(ctx context.Context, user *model.User, ttl time.Duration) error
	Delete(ctx context.Context, id uuid.UUID) error
}
