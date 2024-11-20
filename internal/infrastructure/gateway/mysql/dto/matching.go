package dto

import (
	"errors"

	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/mysql/entity"
	"github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"
)

func ToMatchingModel(entity *entity.MatchingEntity) (*model.Matching, error) {
	if !uuid.IsValidUUIDv7(entity.ID) {
		return nil, errors.New("invalid UUIDv7 format")
	}

	return &model.Matching{
		ID:        uuid.MustParse(entity.ID),
		MeID:      uuid.MustParse(entity.MeID),
		PartnerID: uuid.MustParse(entity.PartnerID),
		Status:    model.MatchingStatus(entity.Status),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}, nil
}

func ToMatchingModels(entities []*entity.MatchingEntity) ([]*model.Matching, error) {
	var users []*model.Matching
	for _, entity := range entities {
		user, err := ToMatchingModel(entity)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func ToMatchingEntity(model *model.Matching) *entity.MatchingEntity {
	return &entity.MatchingEntity{
		ID:        model.ID.String(),
		MeID:      model.MeID.String(),
		PartnerID: model.PartnerID.String(),
		Status:    string(model.Status),
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func ToMatchingEntities(models []*model.Matching) []*entity.MatchingEntity {
	var entities []*entity.MatchingEntity
	for _, model := range models {
		entities = append(entities, ToMatchingEntity(model))
	}
	return entities
}
