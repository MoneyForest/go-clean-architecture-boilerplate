package service

import (
	"context"

	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/repository"
	"github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"
)

type MatchingDomainService struct {
	userRepo  repository.UserRepository
	matchRepo repository.MatchRepository
}

func NewMatchingDomainService(ur repository.UserRepository, mr repository.MatchRepository) *MatchingDomainService {
	return &MatchingDomainService{
		userRepo:  ur,
		matchRepo: mr,
	}
}

func (s *MatchingDomainService) CreateMatch(ctx context.Context, meID, partnerID uuid.UUID) (*model.Match, error) {
	_, err := s.userRepo.Get(ctx, meID)
	if err != nil {
		return nil, err
	}
	_, err = s.userRepo.Get(ctx, partnerID)
	if err != nil {
		return nil, err
	}

	match := model.NewMatch(model.InputMatchParams{
		MeID:      meID,
		PartnerID: partnerID,
		Status:    "pending",
	})

	return match, nil
}
