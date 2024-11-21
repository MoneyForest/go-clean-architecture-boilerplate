package input

import (
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
)

type CreateMatchingInput struct {
	MeID      uuid.UUID `json:"user1_id"`
	PartnerID uuid.UUID `json:"user2_id"`
}

type GetMatchingInput struct {
	ID uuid.UUID `json:"id"`
}

type ListMatchingInput struct {
	UserID uuid.UUID `json:"user_id"`
	Limit  int       `json:"limit"`
	Offset int       `json:"offset"`
}

type UpdateMatchingInput struct {
	ID     uuid.UUID            `json:"id"`
	Status model.MatchingStatus `json:"status"`
}

type DeleteMatchingInput struct {
	ID uuid.UUID `json:"id"`
}
