package model

import (
	"errors"
	"time"

	"github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"
)

var (
	ErrMatchMeOrPartnerIDIsRequired = errors.New("match me or partner id is required")
	ErrMatchStatusIsRequired        = errors.New("match status is required")
	ErrMatchStatusIsInvalid         = errors.New("match status is invalid")
)

type MatchStatus string

const (
	MatchStatusPending  MatchStatus = "pending"
	MatchStatusAccepted MatchStatus = "accepted"
)

var MatchStatuses = map[MatchStatus]struct{}{
	MatchStatusPending:  {},
	MatchStatusAccepted: {},
}

type Match struct {
	ID        uuid.UUID
	MeID      uuid.UUID
	PartnerID uuid.UUID
	Status    MatchStatus
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
		Status:    MatchStatus(params.Status),
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
	if _, ok := MatchStatuses[m.Status]; !ok {
		return ErrMatchStatusIsInvalid
	}
	return nil
}
