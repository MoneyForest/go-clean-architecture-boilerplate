package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/gateway/mysql/sqlc"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/gateway/mysql/transaction"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
)

type MatchingMySQLRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func NewMatchingMySQLRepository(db *sql.DB) *MatchingMySQLRepository {
	return &MatchingMySQLRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *MatchingMySQLRepository) Save(ctx context.Context, matching *model.Matching) (*model.Matching, error) {
	q := transaction.GetQueries(ctx, r.queries)
	exists, err := q.ExistsMatching(ctx, matching.ID.String())
	if err != nil {
		return nil, err
	}

	if exists {
		_, err = q.UpdateMatching(ctx, sqlc.UpdateMatchingParams{
			Status:    string(matching.Status),
			UpdatedAt: time.Now(),
			ID:        matching.ID.String(),
		})
	} else {
		_, err = q.CreateMatching(ctx, sqlc.CreateMatchingParams{
			ID:        matching.ID.String(),
			MeID:      matching.MeID.String(),
			PartnerID: matching.PartnerID.String(),
			Status:    string(matching.Status),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	if err != nil {
		return nil, err
	}
	return matching, nil
}

func (r *MatchingMySQLRepository) FindAllByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Matching, error) {
	q := transaction.GetQueries(ctx, r.queries)
	matchings, err := q.ListMatchingsByUser(ctx, sqlc.ListMatchingsByUserParams{
		MeID:      userID.String(),
		PartnerID: userID.String(),
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*model.Matching, len(matchings))
	for i, m := range matchings {
		result[i] = &model.Matching{
			ID:        uuid.MustParse(m.ID),
			MeID:      uuid.MustParse(m.MeID),
			PartnerID: uuid.MustParse(m.PartnerID),
			Status:    model.MatchingStatus(m.Status),
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
	}
	return result, nil
}

func (r *MatchingMySQLRepository) FindById(ctx context.Context, id uuid.UUID) (*model.Matching, error) {
	q := transaction.GetQueries(ctx, r.queries)
	matching, err := q.GetMatching(ctx, id.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &model.Matching{
		ID:        uuid.MustParse(matching.ID),
		MeID:      uuid.MustParse(matching.MeID),
		PartnerID: uuid.MustParse(matching.PartnerID),
		Status:    model.MatchingStatus(matching.Status),
		CreatedAt: matching.CreatedAt,
		UpdatedAt: matching.UpdatedAt,
	}, nil
}

func (r *MatchingMySQLRepository) FindByParticipants(ctx context.Context, meID, partnerID uuid.UUID) (*model.Matching, error) {
	q := transaction.GetQueries(ctx, r.queries)
	matching, err := q.GetMatchingByParticipants(ctx, sqlc.GetMatchingByParticipantsParams{
		MeID:      meID.String(),
		PartnerID: partnerID.String(),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &model.Matching{
		ID:        uuid.MustParse(matching.ID),
		MeID:      uuid.MustParse(matching.MeID),
		PartnerID: uuid.MustParse(matching.PartnerID),
		Status:    model.MatchingStatus(matching.Status),
		CreatedAt: matching.CreatedAt,
		UpdatedAt: matching.UpdatedAt,
	}, nil
}

func (r *MatchingMySQLRepository) Remove(ctx context.Context, id uuid.UUID) (*uuid.UUID, error) {
	q := transaction.GetQueries(ctx, r.queries)
	err := q.DeleteMatching(ctx, id.String())
	if err != nil {
		return nil, err
	}

	return &id, nil
}
