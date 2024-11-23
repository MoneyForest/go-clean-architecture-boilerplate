package service

import (
	"context"
	"errors"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
)

var (
	ErrMatchingMeAndPartnerAreSameUser = errors.New("me and partner are the same user")
)

type MatchingDomainService struct{}

func (s *MatchingDomainService) Validate(ctx context.Context, me, partner *model.User) error {
	// Write the business logic for the match.
	// For example, a match for an unsubscribed user or a user who violates the Terms of Service is invalid.
	// It can only be validated by a call to the domain model.
	if err := me.Validate(); err != nil {
		return err
	}
	if err := partner.Validate(); err != nil {
		return err
	}
	if me.ID == partner.ID {
		return ErrMatchingMeAndPartnerAreSameUser
	}
	return nil
}
