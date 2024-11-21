package repository

import (
	"context"
	"database/sql"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/gateway/mysql/dto"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/gateway/mysql/entity"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
)

type UserMySQLRepository struct {
	db *sql.DB
}

func NewUserMySQLRepository(db *sql.DB) UserMySQLRepository {
	return UserMySQLRepository{db: db}
}

func (r UserMySQLRepository) Save(ctx context.Context, user *model.User) (*model.User, error) {
	if exists, _ := r.exists(ctx, user.ID); exists {
		return r.update(ctx, user)
	}
	return r.create(ctx, user)
}

func (r UserMySQLRepository) exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM user WHERE id = ?)`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

func (r UserMySQLRepository) create(ctx context.Context, user *model.User) (*model.User, error) {
	query := `INSERT INTO user (id, email, created_at, updated_at) VALUES (?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r UserMySQLRepository) update(ctx context.Context, user *model.User) (*model.User, error) {
	query := `UPDATE user SET email = ?, updated_at = NOW() WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, user.Email, user.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r UserMySQLRepository) FindAll(ctx context.Context, limit, offset int) ([]*model.User, error) {
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

func (r UserMySQLRepository) FindById(ctx context.Context, id uuid.UUID) (*model.User, error) {
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

func (r UserMySQLRepository) Remove(ctx context.Context, id uuid.UUID) (*uuid.UUID, error) {
	query := `DELETE FROM user WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}
