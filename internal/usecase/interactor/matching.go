package interactor

import (
	"context"
	"time"

	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/repository"
	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/service"
	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/transaction"
	"github.com/MoneyForest/go-clean-boilerplate/internal/usecase/port/input"
	"github.com/MoneyForest/go-clean-boilerplate/internal/usecase/port/output"
	"github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"
)

type MatchingInteractor interface {
	Create(ctx context.Context, input *input.CreateMatchingInput) (*output.CreateMatchingOutput, error)
	Get(ctx context.Context, input *input.GetMatchingInput) (*output.GetMatchingOutput, error)
	List(ctx context.Context, input *input.ListMatchingInput) (*output.ListMatchingOutput, error)
	Update(ctx context.Context, input *input.UpdateMatchingInput) (*output.UpdateMatchingOutput, error)
	Delete(ctx context.Context, input *input.DeleteMatchingInput) (*output.DeleteMatchingOutput, error)
}

type matchingInteractor struct {
	txManager transaction.Manager
	repo      repository.MatchingRepository
	service   *service.MatchingDomainService
}

func NewMatchingInteractor(txManager transaction.Manager, repo repository.MatchingRepository, ds *service.MatchingDomainService) MatchingInteractor {
	return &matchingInteractor{
		txManager: txManager,
		repo:      repo,
		service:   ds,
	}
}

func (i *matchingInteractor) Create(ctx context.Context, input *input.CreateMatchingInput) (*output.CreateMatchingOutput, error) {
	matching, err := i.service.CreateMatching(ctx, input.MeID, input.PartnerID)
	if err != nil {
		return nil, err
	}
	var createdMatching *model.Matching
	err = i.txManager.DoInTx(ctx, func(ctx context.Context) error {
		var err error
		createdMatching, err = i.repo.Create(ctx, matching)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &output.CreateMatchingOutput{Matching: createdMatching}, nil
}

func (i *matchingInteractor) Get(ctx context.Context, input *input.GetMatchingInput) (*output.GetMatchingOutput, error) {
	matching, err := i.repo.Get(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return &output.GetMatchingOutput{Matching: matching}, nil
}

func (i *matchingInteractor) List(ctx context.Context, input *input.ListMatchingInput) (*output.ListMatchingOutput, error) {
	matchinges, err := i.repo.List(ctx, input.UserID, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}
	return &output.ListMatchingOutput{Matchinges: matchinges}, nil
}

func (i *matchingInteractor) Update(ctx context.Context, input *input.UpdateMatchingInput) (*output.UpdateMatchingOutput, error) {
	var updatedMatching *model.Matching
	err := i.txManager.DoInTx(ctx, func(ctx context.Context) error {
		matching, err := i.repo.Get(ctx, input.ID)
		if err != nil {
			return err
		}

		matching.Status = model.MatchingStatus(input.Status)
		matching.UpdatedAt = time.Now()

		updatedMatching, err = i.repo.Update(ctx, matching)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &output.UpdateMatchingOutput{Matching: updatedMatching}, nil
}

func (i *matchingInteractor) Delete(ctx context.Context, input *input.DeleteMatchingInput) (*output.DeleteMatchingOutput, error) {
	var deletedID *uuid.UUID
	err := i.txManager.DoInTx(ctx, func(ctx context.Context) error {
		var err error
		deletedID, err = i.repo.Delete(ctx, input.ID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &output.DeleteMatchingOutput{ID: deletedID}, nil
}
