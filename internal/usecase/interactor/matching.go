package interactor

import (
	"context"
	"log"
	"time"

	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/repository"
	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/service"
	"github.com/MoneyForest/go-clean-boilerplate/internal/usecase/port/input"
	"github.com/MoneyForest/go-clean-boilerplate/internal/usecase/port/output"
)

type MatchingInteractor interface {
	Create(ctx context.Context, input *input.CreateMatchingInput) (*output.CreateMatchingOutput, error)
	Get(ctx context.Context, input *input.GetMatchingInput) (*output.GetMatchingOutput, error)
	List(ctx context.Context, input *input.ListMatchingInput) (*output.ListMatchingOutput, error)
	Update(ctx context.Context, input *input.UpdateMatchingInput) (*output.UpdateMatchingOutput, error)
	Delete(ctx context.Context, input *input.DeleteMatchingInput) (*output.DeleteMatchingOutput, error)
}

type matchingInteractor struct {
	repo    repository.MatchingRepository
	service *service.MatchingDomainService
}

func NewMatchingInteractor(repo repository.MatchingRepository, ds *service.MatchingDomainService) MatchingInteractor {
	return &matchingInteractor{
		repo:    repo,
		service: ds,
	}
}

func (i *matchingInteractor) Create(ctx context.Context, input *input.CreateMatchingInput) (*output.CreateMatchingOutput, error) {
	tx, err := i.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	matching, err := i.service.CreateMatching(ctx, input.MeID, input.PartnerID)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("failed to rollback transaction: %v\n", err)
		}
	}()

	createdMatching, err := i.repo.CreateTx(ctx, tx, matching)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
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
	matching, err := i.repo.Get(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	matching.Status = model.MatchingStatus(input.Status)
	matching.UpdatedAt = time.Now()

	tx, err := i.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("failed to rollback transaction: %v\n", err)
		}
	}()

	updatedMatching, err := i.repo.UpdateTx(ctx, tx, matching)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &output.UpdateMatchingOutput{Matching: updatedMatching}, nil
}

func (i *matchingInteractor) Delete(ctx context.Context, input *input.DeleteMatchingInput) (*output.DeleteMatchingOutput, error) {
	tx, err := i.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("failed to rollback transaction: %v\n", err)
		}
	}()

	deletedID, err := i.repo.DeleteTx(ctx, tx, input.ID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &output.DeleteMatchingOutput{ID: deletedID}, nil
}
