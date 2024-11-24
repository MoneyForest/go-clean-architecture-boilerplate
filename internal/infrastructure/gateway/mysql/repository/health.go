package repository

import (
	"context"
	"database/sql"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/gateway/mysql/sqlc"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/gateway/mysql/transaction"
)

type HealthMySQLRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func NewHealthMySQLRepository(db *sql.DB) *HealthMySQLRepository {
	return &HealthMySQLRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *HealthMySQLRepository) Ping(ctx context.Context) error {
	q := transaction.GetQueries(ctx, r.queries)
	_, err := q.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}
