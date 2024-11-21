package dto

import (
	"encoding/json"
	"errors"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/gateway/redis/entity"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
)

func ToUserModel(entity *entity.UserEntity) (*model.User, error) {
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

func ToUserModels(entities []*entity.UserEntity) ([]*model.User, error) {
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

func ToUserEntity(model *model.User) *entity.UserEntity {
	return &entity.UserEntity{
		ID:        model.ID.String(),
		Email:     model.Email,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func ToUserEntities(models []*model.User) []*entity.UserEntity {
	var entities []*entity.UserEntity
	for _, model := range models {
		entities = append(entities, ToUserEntity(model))
	}
	return entities
}

func ToJSON(entity *entity.UserEntity) (string, error) {
	bytes, err := json.Marshal(entity)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func FromJSON(data string) (*entity.UserEntity, error) {
	var entity entity.UserEntity
	if err := json.Unmarshal([]byte(data), &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}
