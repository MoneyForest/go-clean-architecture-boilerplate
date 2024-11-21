package repository

import (
	"context"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
)

type MatchingRepository interface {
	Create(ctx context.Context, matching *model.Matching) (*model.Matching, error)
	Get(ctx context.Context, id uuid.UUID) (*model.Matching, error)
	List(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Matching, error)
	Update(ctx context.Context, matching *model.Matching) (*model.Matching, error)
	Delete(ctx context.Context, id uuid.UUID) (*uuid.UUID, error)
}
