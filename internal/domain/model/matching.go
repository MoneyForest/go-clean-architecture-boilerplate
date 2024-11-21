package model

import (
	"errors"
	"time"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
)

var (
	ErrMatchingMeOrPartnerIDIsRequired = errors.New("matching me or partner id is required")
	ErrMatchingStatusIsRequired        = errors.New("matching status is required")
	ErrMatchingStatusIsInvalid         = errors.New("matching status is invalid")
)

type MatchingStatus string

const (
	MatchingStatusPending  MatchingStatus = "pending"
	MatchingStatusAccepted MatchingStatus = "accepted"
	MatchingStatusRejected MatchingStatus = "rejected"
)

var MatchingStatuses = map[MatchingStatus]struct{}{
	MatchingStatusPending:  {},
	MatchingStatusAccepted: {},
}

type Matching struct {
	ID        uuid.UUID
	MeID      uuid.UUID
	PartnerID uuid.UUID
	Status    MatchingStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type InputMatchingParams struct {
	ID        uuid.UUID
	MeID      uuid.UUID
	PartnerID uuid.UUID
	Status    string
}

func NewMatching(params InputMatchingParams) *Matching {
	if params.ID == uuid.Nil() {
		params.ID = uuid.New()
	}
	return &Matching{
		ID:        params.ID,
		MeID:      params.MeID,
		PartnerID: params.PartnerID,
		Status:    MatchingStatus(params.Status),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (m *Matching) Validate() error {
	if m.MeID == uuid.Nil() || m.PartnerID == uuid.Nil() {
		return ErrMatchingMeOrPartnerIDIsRequired
	}
	if m.Status == "" {
		return ErrMatchingStatusIsRequired
	}
	if _, ok := MatchingStatuses[m.Status]; !ok {
		return ErrMatchingStatusIsInvalid
	}
	return nil
}

func (m *Matching) Accept() {
	m.Status = MatchingStatusAccepted
}

func (m *Matching) Reject() {
	m.Status = MatchingStatusRejected
}
