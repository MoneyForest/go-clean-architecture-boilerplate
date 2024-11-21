package port

import (
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
)

type CreateMatchingInput struct {
	MeID      uuid.UUID `json:"user1_id"`
	PartnerID uuid.UUID `json:"user2_id"`
}

type CreateMatchingOutput struct {
	Matching *model.Matching `json:"matching"`
}

type AcceptMatchingInput struct {
	MeID      uuid.UUID `json:"me_id"`
	PartnerID uuid.UUID `json:"partner_id"`
}

type AcceptMatchingOutput struct {
	Matching *model.Matching `json:"matching"`
}

type RejectMatchingInput struct {
	MeID      uuid.UUID `json:"me_id"`
	PartnerID uuid.UUID `json:"partner_id"`
}

type RejectMatchingOutput struct {
	Matching *model.Matching `json:"matching"`
}

type ListMatchingByMeIDInput struct {
	MeID   uuid.UUID `json:"me_id"`
	Limit  int       `json:"limit"`
	Offset int       `json:"offset"`
}

type ListMatchingByMeIDOutput struct {
	Matchings []*model.Matching `json:"matchings"`
}
