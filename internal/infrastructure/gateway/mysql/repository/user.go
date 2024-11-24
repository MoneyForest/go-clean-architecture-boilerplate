package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/gateway/mysql/sqlc"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/gateway/mysql/transaction"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
)

type UserMySQLRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func NewUserMySQLRepository(db *sql.DB) *UserMySQLRepository {
	return &UserMySQLRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *UserMySQLRepository) Save(ctx context.Context, user *model.User) (*model.User, error) {
	q := transaction.GetQueries(ctx, r.queries)
	exists, err := q.ExistsUser(ctx, user.ID.String())
	if err != nil {
		return nil, err
	}

	if exists {
		_, err = q.UpdateUser(ctx, sqlc.UpdateUserParams{
			Email:     user.Email,
			UpdatedAt: time.Now(),
			ID:        user.ID.String(),
		})
	} else {
		_, err = q.CreateUser(ctx, sqlc.CreateUserParams{
			ID:        user.ID.String(),
			Email:     user.Email,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserMySQLRepository) FindAll(ctx context.Context, limit, offset int) ([]*model.User, error) {
	q := transaction.GetQueries(ctx, r.queries)
	users, err := q.ListUsers(ctx, sqlc.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*model.User, len(users))
	for i, user := range users {
		result[i] = &model.User{
			ID:        uuid.MustParse(user.ID),
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}
	return result, nil
}

func (r *UserMySQLRepository) FindById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	q := transaction.GetQueries(ctx, r.queries)
	user, err := q.GetUser(ctx, id.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &model.User{
		ID:        uuid.MustParse(user.ID),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (r *UserMySQLRepository) Remove(ctx context.Context, id uuid.UUID) (*uuid.UUID, error) {
	q := transaction.GetQueries(ctx, r.queries)
	err := q.DeleteUser(ctx, id.String())
	if err != nil {
		return nil, err
	}

	return &id, nil
}
