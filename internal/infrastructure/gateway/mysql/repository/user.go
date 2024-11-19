package repository

import (
	"context"
	"database/sql"

	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/mysql/dto"
	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/mysql/entity"
	"github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"
)

type UserMySQLRepository struct {
	db *sql.DB
}

func NewUserMySQLRepository(db *sql.DB) UserMySQLRepository {
	return UserMySQLRepository{db: db}
}

func (r UserMySQLRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r UserMySQLRepository) CreateTx(ctx context.Context, tx *sql.Tx, user *model.User) (*model.User, error) {
	query := `INSERT INTO user (id, email, created_at, updated_at) VALUES (?, ?, ?, ?)`

	_, err := tx.ExecContext(ctx, query, user.ID, user.Email, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r UserMySQLRepository) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	query := `SELECT id, email, created_at, updated_at FROM user LIMIT ? OFFSET ?`

	var entities []*entity.UserEntity
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entity entity.UserEntity
		if err := rows.Scan(
			&entity.ID,
			&entity.Email,
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

	return dto.ToUserModels(entities)
}

func (r UserMySQLRepository) Get(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `SELECT id, email, created_at, updated_at FROM user WHERE id = ?`

	var entity entity.UserEntity
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&entity.ID,
		&entity.Email,
		&entity.CreatedAt,
		&entity.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return dto.ToUserModel(&entity)
}

func (r UserMySQLRepository) UpdateTx(ctx context.Context, tx *sql.Tx, user *model.User) (*model.User, error) {
	query := `UPDATE user SET email = ?, updated_at = NOW() WHERE id = ?`

	_, err := tx.ExecContext(ctx, query, user.Email, user.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r UserMySQLRepository) DeleteTx(ctx context.Context, tx *sql.Tx, id uuid.UUID) (*uuid.UUID, error) {
	query := `DELETE FROM user WHERE id = ?`

	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}
