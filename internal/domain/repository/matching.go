package repository

import (
	"context"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
)

type MatchingRepository interface {
	Save(ctx context.Context, matching *model.Matching) (*model.Matching, error)
	FindById(ctx context.Context, id uuid.UUID) (*model.Matching, error)
	FindByParticipants(ctx context.Context, meID, partnerID uuid.UUID) (*model.Matching, error)
	FindAllByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Matching, error)
	Remove(ctx context.Context, id uuid.UUID) (*uuid.UUID, error)
}
