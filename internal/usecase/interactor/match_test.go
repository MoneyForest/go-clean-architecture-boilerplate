package interactor

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/service"
	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/mysql/repository"
	"github.com/MoneyForest/go-clean-boilerplate/internal/usecase/port/input"
	"github.com/MoneyForest/go-clean-boilerplate/pkg/testhelper"
	"github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"
)

func SetupTestMatchInteractor(ctx context.Context, gw *testhelper.Gateway) (MatchInteractor, *repository.UserMySQLRepository) {
	matchRepo := repository.NewMatchMySQLRepository(gw.MySQLClient) // ポインタを返すように &を追加
	userRepo := repository.NewUserMySQLRepository(gw.MySQLClient)
	ds := service.NewMatchingDomainService(userRepo, matchRepo)
	return NewMatchInteractor(matchRepo, ds), &userRepo
}

func createTestUser(ctx context.Context, t *testing.T, userRepo *repository.UserMySQLRepository) *model.User {
	ID := uuid.New()
	user := model.NewUser(model.InputUserParams{
		ID:    ID,
		Email: ID.String() + "@example.com",
	})
	tx, err := userRepo.BeginTx(ctx)
	if err != nil {
		t.Fatalf("Failed to begin tx: %v", err)
	}

	createdUser, err := userRepo.CreateTx(ctx, tx, user)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Failed to create test user: %v", err)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		t.Fatalf("Failed to commit tx: %v", err)
	}

	return createdUser
}

func TestMatchInteractor_Create(t *testing.T) {
	ctx := context.Background()
	gw, err := testhelper.Setup(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test: %v", err)
	}
	defer testhelper.Cleanup(ctx, gw)
	matchInteractor, userRepo := SetupTestMatchInteractor(ctx, gw)

	user1 := createTestUser(ctx, t, userRepo)
	user2 := createTestUser(ctx, t, userRepo)

	tests := []struct {
		name    string
		input   *input.CreateMatchInput
		wantErr bool
	}{
		{
			name: "OK",
			input: &input.CreateMatchInput{
				MeID:      user1.ID,
				PartnerID: user2.ID,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchInteractor.Create(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.Match == nil {
				t.Error("Create() got nil match")
				return
			}
			if got.Match.MeID != tt.input.MeID || got.Match.PartnerID != tt.input.PartnerID {
				t.Errorf("Create() got = %v, want meID: %v, partnerID: %v", got.Match, tt.input.MeID, tt.input.PartnerID)
			}
		})
	}
}

func TestMatchInteractor_Get(t *testing.T) {
	ctx := context.Background()
	gw, err := testhelper.Setup(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test: %v", err)
	}
	defer testhelper.Cleanup(ctx, gw)
	matchInteractor, userRepo := SetupTestMatchInteractor(ctx, gw)

	user1 := createTestUser(ctx, t, userRepo)
	user2 := createTestUser(ctx, t, userRepo)

	created, err := matchInteractor.Create(ctx, &input.CreateMatchInput{
		MeID:      user1.ID,
		PartnerID: user2.ID,
	})
	if err != nil {
		t.Fatalf("Failed to create test match: %v", err)
	}

	tests := []struct {
		name    string
		input   *input.GetMatchInput
		want    *model.Match
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &input.GetMatchInput{ID: created.Match.ID},
			want:    created.Match,
			wantErr: false,
		},
		{
			name:    "NotFound",
			input:   &input.GetMatchInput{ID: uuid.New()},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchInteractor.Get(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			diff := cmp.Diff(
				got.Match,
				tt.want,
				cmpopts.IgnoreFields(model.Match{}, "ID", "CreatedAt", "UpdatedAt"),
			)
			if diff != "" {
				t.Errorf("Get() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestMatchInteractor_List(t *testing.T) {
	ctx := context.Background()
	gw, err := testhelper.Setup(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test: %v", err)
	}
	defer testhelper.Cleanup(ctx, gw)
	matchInteractor, userRepo := SetupTestMatchInteractor(ctx, gw)

	user1 := createTestUser(ctx, t, userRepo)
	user2 := createTestUser(ctx, t, userRepo)

	for i := 0; i < 3; i++ {
		partner := createTestUser(ctx, t, userRepo)
		_, err := matchInteractor.Create(ctx, &input.CreateMatchInput{
			MeID:      user1.ID,
			PartnerID: partner.ID,
		})
		if err != nil {
			t.Fatalf("Failed to create test match: %v", err)
		}
	}

	tests := []struct {
		name    string
		input   *input.ListMatchInput
		want    int
		wantErr bool
	}{
		{
			name: "OK_AllMatches",
			input: &input.ListMatchInput{
				UserID: user1.ID,
				Limit:  10,
				Offset: 0,
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "OK_LimitedMatches",
			input: &input.ListMatchInput{
				UserID: user1.ID,
				Limit:  2,
				Offset: 0,
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "OK_NoMatches",
			input: &input.ListMatchInput{
				UserID: user2.ID,
				Limit:  10,
				Offset: 0,
			},
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchInteractor.List(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got.Matches) != tt.want {
				t.Errorf("List() got = %v matches, want %v", len(got.Matches), tt.want)
			}
		})
	}
}

func TestMatchInteractor_Update(t *testing.T) {
	ctx := context.Background()
	gw, err := testhelper.Setup(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test: %v", err)
	}
	defer testhelper.Cleanup(ctx, gw)
	matchInteractor, userRepo := SetupTestMatchInteractor(ctx, gw)

	user1 := createTestUser(ctx, t, userRepo)
	user2 := createTestUser(ctx, t, userRepo)

	created, err := matchInteractor.Create(ctx, &input.CreateMatchInput{
		MeID:      user1.ID,
		PartnerID: user2.ID,
	})
	if err != nil {
		t.Fatalf("Failed to create test match: %v", err)
	}

	tests := []struct {
		name    string
		input   *input.UpdateMatchInput
		want    model.MatchStatus
		wantErr bool
	}{
		{
			name: "OK",
			input: &input.UpdateMatchInput{
				ID:     created.Match.ID,
				Status: model.MatchStatusAccepted,
			},
			want:    model.MatchStatusAccepted,
			wantErr: false,
		},
		{
			name: "NG_NotFound",
			input: &input.UpdateMatchInput{
				ID:     uuid.New(),
				Status: model.MatchStatusAccepted,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchInteractor.Update(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.Match.Status != tt.want {
				t.Errorf("Update() got status = %v, want %v", got.Match.Status, tt.want)
			}
		})
	}
}

func TestMatchInteractor_Delete(t *testing.T) {
	ctx := context.Background()
	gw, err := testhelper.Setup(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test: %v", err)
	}
	defer testhelper.Cleanup(ctx, gw)
	matchInteractor, userRepo := SetupTestMatchInteractor(ctx, gw)

	user1 := createTestUser(ctx, t, userRepo)
	user2 := createTestUser(ctx, t, userRepo)

	created, err := matchInteractor.Create(ctx, &input.CreateMatchInput{
		MeID:      user1.ID,
		PartnerID: user2.ID,
	})
	if err != nil {
		t.Fatalf("Failed to create test match: %v", err)
	}

	tests := []struct {
		name    string
		input   *input.DeleteMatchInput
		wantErr bool
	}{
		{
			name: "OK",
			input: &input.DeleteMatchInput{
				ID: created.Match.ID,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchInteractor.Delete(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if *got.ID != tt.input.ID {
				t.Errorf("Delete() got = %v, want %v", got.ID, tt.input.ID)
			}
		})
	}
}
