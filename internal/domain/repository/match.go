package repository

import (
	"context"
	"database/sql"

	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"
)

type MatchRepository interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	CreateTx(ctx context.Context, tx *sql.Tx, match *model.Match) (*model.Match, error)
	Get(ctx context.Context, id uuid.UUID) (*model.Match, error)
	List(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Match, error)
	UpdateTx(ctx context.Context, tx *sql.Tx, match *model.Match) (*model.Match, error)
	DeleteTx(ctx context.Context, tx *sql.Tx, id uuid.UUID) (*uuid.UUID, error)
}
