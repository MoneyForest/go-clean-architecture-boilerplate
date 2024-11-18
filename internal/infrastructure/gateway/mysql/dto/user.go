package dto

import (
	"errors"
	"time"

	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"
)

type UserEntity struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func ToUserModel(entity *UserEntity) (*model.User, error) {
	if !uuid.IsValidUUIDv7(entity.ID) {
		return nil, errors.New("invalid UUIDv7 format")
	}

	return &model.User{
		ID:        uuid.MustParse(entity.ID),
		Email:     entity.Email,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}, nil
}

func ToUserModels(entities []*UserEntity) ([]*model.User, error) {
	var users []*model.User
	for _, entity := range entities {
		user, err := ToUserModel(entity)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func ToUserEntity(model *model.User) *UserEntity {
	return &UserEntity{
		ID:        model.ID.String(),
		Email:     model.Email,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
