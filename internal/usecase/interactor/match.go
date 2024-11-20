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

type MatchInteractor interface {
	Create(ctx context.Context, input *input.CreateMatchInput) (*output.CreateMatchOutput, error)
	Get(ctx context.Context, input *input.GetMatchInput) (*output.GetMatchOutput, error)
	List(ctx context.Context, input *input.ListMatchInput) (*output.ListMatchOutput, error)
	Update(ctx context.Context, input *input.UpdateMatchInput) (*output.UpdateMatchOutput, error)
	Delete(ctx context.Context, input *input.DeleteMatchInput) (*output.DeleteMatchOutput, error)
}

type matchInteractor struct {
	mysql         repository.MatchRepository
	domainService *service.MatchingDomainService
}

func NewMatchInteractor(mysql repository.MatchRepository, ds *service.MatchingDomainService) MatchInteractor {
	return &matchInteractor{
		mysql:         mysql,
		domainService: ds,
	}
}

func (i *matchInteractor) Create(ctx context.Context, input *input.CreateMatchInput) (*output.CreateMatchOutput, error) {
	tx, err := i.mysql.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	match, err := i.domainService.CreateMatch(ctx, input.MeID, input.PartnerID)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("failed to rollback transaction: %v\n", err)
		}
	}()

	createdMatch, err := i.mysql.CreateTx(ctx, tx, match)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &output.CreateMatchOutput{Match: createdMatch}, nil
}

func (i *matchInteractor) Get(ctx context.Context, input *input.GetMatchInput) (*output.GetMatchOutput, error) {
	match, err := i.mysql.Get(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return &output.GetMatchOutput{Match: match}, nil
}

func (i *matchInteractor) List(ctx context.Context, input *input.ListMatchInput) (*output.ListMatchOutput, error) {
	matches, err := i.mysql.List(ctx, input.UserID, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}
	return &output.ListMatchOutput{Matches: matches}, nil
}

func (i *matchInteractor) Update(ctx context.Context, input *input.UpdateMatchInput) (*output.UpdateMatchOutput, error) {
	match, err := i.mysql.Get(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	match.Status = model.MatchStatus(input.Status)
	match.UpdatedAt = time.Now()

	tx, err := i.mysql.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("failed to rollback transaction: %v\n", err)
		}
	}()

	updatedMatch, err := i.mysql.UpdateTx(ctx, tx, match)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &output.UpdateMatchOutput{Match: updatedMatch}, nil
}

func (i *matchInteractor) Delete(ctx context.Context, input *input.DeleteMatchInput) (*output.DeleteMatchOutput, error) {
	tx, err := i.mysql.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("failed to rollback transaction: %v\n", err)
		}
	}()

	deletedID, err := i.mysql.DeleteTx(ctx, tx, input.ID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &output.DeleteMatchOutput{ID: deletedID}, nil
}
