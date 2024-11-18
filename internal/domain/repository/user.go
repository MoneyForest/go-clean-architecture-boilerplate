package repository

import (
	"context"

	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Get(ctx context.Context, id uuid.UUID) (*model.User, error)
	List(ctx context.Context, limit, offset int) ([]*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	Dequeue(ctx context.Context, queueURL string) (*model.User, error)
}
