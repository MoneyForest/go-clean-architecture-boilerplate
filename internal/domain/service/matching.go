package service

import (
	"context"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/repository"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
)

type MatchingDomainService struct {
	userRepo     repository.UserRepository
	matchingRepo repository.MatchingRepository
}

func NewMatchingDomainService(ur repository.UserRepository, mr repository.MatchingRepository) *MatchingDomainService {
	return &MatchingDomainService{
		userRepo:     ur,
		matchingRepo: mr,
	}
}

func (s *MatchingDomainService) CreateMatching(ctx context.Context, meID, partnerID uuid.UUID) (*model.Matching, error) {
	me, err := s.userRepo.FindById(ctx, meID)
	if err != nil {
		return nil, err
	}
	partner, err := s.userRepo.FindById(ctx, partnerID)
	if err != nil {
		return nil, err
	}

	matching := model.NewMatching(model.InputMatchingParams{
		MeID:      me.ID,
		PartnerID: partner.ID,
		Status:    "pending",
	})
	createdMatching, err := s.matchingRepo.Save(ctx, matching)
	if err != nil {
		return nil, err
	}

	return createdMatching, nil
}
