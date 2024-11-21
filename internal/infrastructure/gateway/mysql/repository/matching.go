package repository

import (
	"context"
	"database/sql"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/gateway/mysql/dto"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/gateway/mysql/entity"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
)

type MatchingMySQLRepository struct {
	db *sql.DB
}

func NewMatchingMySQLRepository(db *sql.DB) MatchingMySQLRepository {
	return MatchingMySQLRepository{db: db}
}

func (r MatchingMySQLRepository) Save(ctx context.Context, matching *model.Matching) (*model.Matching, error) {
	if exists, _ := r.exists(ctx, matching.ID); exists {
		return r.update(ctx, matching)
	}
	return r.create(ctx, matching)
}

func (r MatchingMySQLRepository) exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM matching WHERE id = ?)`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

func (r MatchingMySQLRepository) create(ctx context.Context, matching *model.Matching) (*model.Matching, error) {
	query := `INSERT INTO matching (id, me_id, partner_id, status, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		matching.ID, matching.MeID, matching.PartnerID,
		matching.Status, matching.CreatedAt, matching.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return matching, nil
}

func (r MatchingMySQLRepository) update(ctx context.Context, matching *model.Matching) (*model.Matching, error) {
	query := `UPDATE matching
              SET status = ?, updated_at = NOW()
              WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, matching.Status, matching.ID)
	if err != nil {
		return nil, err
	}

	return matching, nil
}

func (r MatchingMySQLRepository) FindAllByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Matching, error) {
	query := `SELECT id, me_id, partner_id, status, created_at, updated_at
              FROM matching
              WHERE me_id = ? OR partner_id = ?
              LIMIT ? OFFSET ?`

	var entities []*entity.MatchingEntity
	rows, err := r.db.QueryContext(ctx, query, userID, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entity entity.MatchingEntity
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

	return dto.ToMatchingModels(entities)
}

func (r MatchingMySQLRepository) FindById(ctx context.Context, id uuid.UUID) (*model.Matching, error) {
	query := `SELECT id, me_id, partner_id, status, created_at, updated_at
              FROM matching
              WHERE id = ?`

	var entity entity.MatchingEntity
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

	return dto.ToMatchingModel(&entity)
}

func (r MatchingMySQLRepository) FindByParticipants(ctx context.Context, meID, partnerID uuid.UUID) (*model.Matching, error) {
	query := `SELECT id, me_id, partner_id, status, created_at, updated_at
              FROM matching
              WHERE me_id = ? AND partner_id = ?`

	var entity entity.MatchingEntity
	err := r.db.QueryRowContext(ctx, query, meID, partnerID).Scan(
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

	return dto.ToMatchingModel(&entity)
}

func (r MatchingMySQLRepository) Remove(ctx context.Context, id uuid.UUID) (*uuid.UUID, error) {
	query := `DELETE FROM matching WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}
