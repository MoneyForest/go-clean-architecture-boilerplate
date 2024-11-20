package dto

import (
	"errors"

	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/mysql/entity"
	"github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"
)

func ToMatchModel(entity *entity.MatchEntity) (*model.Match, error) {
	if !uuid.IsValidUUIDv7(entity.ID) {
		return nil, errors.New("invalid UUIDv7 format")
	}

	return &model.Match{
		ID:        uuid.MustParse(entity.ID),
		MeID:      uuid.MustParse(entity.MeID),
		PartnerID: uuid.MustParse(entity.PartnerID),
		Status:    model.MatchStatus(entity.Status),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}, nil
}

func ToMatchModels(entities []*entity.MatchEntity) ([]*model.Match, error) {
	var users []*model.Match
	for _, entity := range entities {
		user, err := ToMatchModel(entity)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func ToMatchEntity(model *model.Match) *entity.MatchEntity {
	return &entity.MatchEntity{
		ID:        model.ID.String(),
		MeID:      model.MeID.String(),
		PartnerID: model.PartnerID.String(),
		Status:    string(model.Status),
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func ToMatchEntities(models []*model.Match) []*entity.MatchEntity {
	var entities []*entity.MatchEntity
	for _, model := range models {
		entities = append(entities, ToMatchEntity(model))
	}
	return entities
}
