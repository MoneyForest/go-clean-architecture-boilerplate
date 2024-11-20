package repository

import (
	"context"
	"database/sql"

	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"
)

type MatchingRepository interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	CreateTx(ctx context.Context, tx *sql.Tx, matching *model.Matching) (*model.Matching, error)
	Get(ctx context.Context, id uuid.UUID) (*model.Matching, error)
	List(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Matching, error)
	UpdateTx(ctx context.Context, tx *sql.Tx, matching *model.Matching) (*model.Matching, error)
	DeleteTx(ctx context.Context, tx *sql.Tx, id uuid.UUID) (*uuid.UUID, error)
}
