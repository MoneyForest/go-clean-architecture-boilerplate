package input

import "github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"

type CreateMatchInput struct {
	MeID      uuid.UUID `json:"user1_id"`
	PartnerID uuid.UUID `json:"user2_id"`
}

type GetMatchInput struct {
	ID uuid.UUID `json:"id"`
}

type ListMatchInput struct {
	UserID uuid.UUID `json:"user_id"`
	Limit  int       `json:"limit"`
	Offset int       `json:"offset"`
}

type UpdateMatchInput struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

type DeleteMatchInput struct {
	ID uuid.UUID `json:"id"`
}
