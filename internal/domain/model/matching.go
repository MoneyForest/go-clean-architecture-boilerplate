package model

import (
	"errors"
	"time"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
	"github.com/go-playground/validator/v10"
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
	ID        uuid.UUID      `validate:"required"`
	MeID      uuid.UUID      `validate:"required"`
	PartnerID uuid.UUID      `validate:"required"`
	Status    MatchingStatus `validate:"required,matching_status"`
	CreatedAt time.Time      `validate:"required"`
	UpdatedAt time.Time      `validate:"required"`
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
	validate := validator.New()
	if err := validate.RegisterValidation("matching_status", validateMatchingStatus); err != nil {
		return err
	}
	if err := validate.Struct(m); err != nil {
		return err
	}
	return nil
}

func validateMatchingStatus(fl validator.FieldLevel) bool {
	status, ok := fl.Field().Interface().(MatchingStatus)
	if !ok {
		return false
	}
	_, exists := MatchingStatuses[status]
	return exists
}

func (m *Matching) Accept() {
	m.Status = MatchingStatusAccepted
}

func (m *Matching) Reject() {
	m.Status = MatchingStatusRejected
}
