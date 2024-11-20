package repository

import (
	"context"
	"database/sql"

	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/mysql/dto"
	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/mysql/entity"
	"github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"
)

type MatchMySQLRepository struct {
	db *sql.DB
}

func NewMatchMySQLRepository(db *sql.DB) MatchMySQLRepository {
	return MatchMySQLRepository{db: db}
}

func (r MatchMySQLRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r MatchMySQLRepository) CreateTx(ctx context.Context, tx *sql.Tx, match *model.Match) (*model.Match, error) {
	query := `INSERT INTO matching (id, me_id, partner_id, status, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?)`

	_, err := tx.ExecContext(ctx, query,
		match.ID, match.MeID, match.PartnerID,
		match.Status, match.CreatedAt, match.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return match, nil
}

func (r MatchMySQLRepository) List(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Match, error) {
	query := `SELECT id, me_id, partner_id, status, created_at, updated_at
              FROM matching
              WHERE me_id = ? OR partner_id = ?
              LIMIT ? OFFSET ?`

	var entities []*entity.MatchEntity
	rows, err := r.db.QueryContext(ctx, query, userID, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entity entity.MatchEntity
		if err := rows.Scan(
			&entity.ID,
			&entity.MeID,
			&entity.PartnerID,
			&entity.Status,
			&entity.CreatedAt,
			&entity.UpdatedAt,
		); err != nil {
			return nil, err
		}
		entities = append(entities, &entity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dto.ToMatchModels(entities)
}

func (r MatchMySQLRepository) Get(ctx context.Context, id uuid.UUID) (*model.Match, error) {
	query := `SELECT id, me_id, partner_id, status, created_at, updated_at
              FROM matching
              WHERE id = ?`

	var entity entity.MatchEntity
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&entity.ID,
		&entity.MeID,
		&entity.PartnerID,
		&entity.Status,
		&entity.CreatedAt,
		&entity.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return dto.ToMatchModel(&entity)
}

func (r MatchMySQLRepository) UpdateTx(ctx context.Context, tx *sql.Tx, match *model.Match) (*model.Match, error) {
	query := `UPDATE matching
              SET status = ?, updated_at = NOW()
              WHERE id = ?`

	_, err := tx.ExecContext(ctx, query, match.Status, match.ID)
	if err != nil {
		return nil, err
	}

	return match, nil
}

func (r MatchMySQLRepository) DeleteTx(ctx context.Context, tx *sql.Tx, id uuid.UUID) (*uuid.UUID, error) {
	query := `DELETE FROM matching WHERE id = ?`

	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}
