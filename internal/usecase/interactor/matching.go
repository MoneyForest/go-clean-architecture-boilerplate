package interactor

import (
	"context"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/repository"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/service"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/usecase/port"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/usecase/transaction"
)

type MatchingInteractor struct {
	txManager transaction.Manager
	repo      repository.MatchingRepository
	userRepo  repository.UserRepository
	service   *service.MatchingDomainService
}

func NewMatchingInteractor(txManager transaction.Manager, repo repository.MatchingRepository, userRepo repository.UserRepository, service *service.MatchingDomainService) MatchingInteractor {
	return MatchingInteractor{
		txManager: txManager,
		repo:      repo,
		userRepo:  userRepo,
		service:   service,
	}
}

func (i MatchingInteractor) Create(ctx context.Context, input *port.CreateMatchingInput) (*port.CreateMatchingOutput, error) {
	var createdMatching *model.Matching
	err := i.txManager.Do(ctx, func(ctx context.Context) error {
		me, err := i.userRepo.FindById(ctx, input.MeID)
		if err != nil {
			return err
		}
		partner, err := i.userRepo.FindById(ctx, input.PartnerID)
		if err != nil {
			return err
		}
		if err := i.service.ValidateMatching(ctx, me, partner); err != nil {
			return err
		}
		createdMatching = model.NewMatching(model.InputMatchingParams{
			MeID:      input.MeID,
			PartnerID: input.PartnerID,
		})
		createdMatching, err = i.repo.Save(ctx, createdMatching)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &port.CreateMatchingOutput{Matching: createdMatching}, nil
}

func (i MatchingInteractor) Accept(ctx context.Context, input *port.AcceptMatchingInput) (*port.AcceptMatchingOutput, error) {
	matching, err := i.repo.FindByParticipants(ctx, input.MeID, input.PartnerID)
	if err != nil {
		return nil, err
	}
	if err := matching.Accept(); err != nil {
		return nil, err
	}

	var updatedMatching *model.Matching
	err = i.txManager.Do(ctx, func(ctx context.Context) error {
		updatedMatching, err = i.repo.Save(ctx, matching)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &port.AcceptMatchingOutput{Matching: updatedMatching}, nil
}

func (i MatchingInteractor) Reject(ctx context.Context, input *port.RejectMatchingInput) (*port.RejectMatchingOutput, error) {
	matching, err := i.repo.FindByParticipants(ctx, input.MeID, input.PartnerID)
	if err != nil {
		return nil, err
	}
	if err := matching.Reject(); err != nil {
		return nil, err
	}

	var updatedMatching *model.Matching
	err = i.txManager.Do(ctx, func(ctx context.Context) error {
		updatedMatching, err = i.repo.Save(ctx, matching)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &port.RejectMatchingOutput{Matching: updatedMatching}, nil
}

func (i MatchingInteractor) ListByMeID(ctx context.Context, input *port.ListMatchingByMeIDInput) (*port.ListMatchingByMeIDOutput, error) {
	matchings, err := i.repo.FindAllByUser(ctx, input.MeID, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}
	return &port.ListMatchingByMeIDOutput{Matchings: matchings}, nil
}
