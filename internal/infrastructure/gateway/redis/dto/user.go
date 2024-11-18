package dto

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"
)

type RedisUserEntity struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToRedisUserModel(entity *RedisUserEntity) (*model.User, error) {
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

func ToRedisUserModels(entities []*RedisUserEntity) ([]*model.User, error) {
	var users []*model.User
	for _, entity := range entities {
		user, err := ToRedisUserModel(entity)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func ToRedisUserEntity(model *model.User) *RedisUserEntity {
	return &RedisUserEntity{
		ID:        model.ID.String(),
		Email:     model.Email,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func (e *RedisUserEntity) ToJSON() (string, error) {
	bytes, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func FromJSON(data string) (*RedisUserEntity, error) {
	var entity RedisUserEntity
	if err := json.Unmarshal([]byte(data), &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}
