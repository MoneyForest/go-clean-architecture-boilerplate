package model

import (
	"errors"
	"time"

	"github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"
)

var (
	ErrMatchMeOrPartnerIDIsRequired = errors.New("match me or partner id is required")
	ErrMatchStatusIsRequired        = errors.New("match status is required")
)

type Match struct {
	ID        uuid.UUID
	MeID      uuid.UUID
	PartnerID uuid.UUID
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type InputMatchParams struct {
	ID        uuid.UUID
	MeID      uuid.UUID
	PartnerID uuid.UUID
	Status    string
}

func NewMatch(params InputMatchParams) *Match {
	if params.ID == uuid.Nil() {
		params.ID = uuid.New()
	}
	return &Match{
		ID:        params.ID,
		MeID:      params.MeID,
		PartnerID: params.PartnerID,
		Status:    params.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (m *Match) Validate() error {
	if m.MeID == uuid.Nil() || m.PartnerID == uuid.Nil() {
		return ErrMatchMeOrPartnerIDIsRequired
	}
	if m.Status == "" {
		return ErrMatchStatusIsRequired
	}
	return nil
}
