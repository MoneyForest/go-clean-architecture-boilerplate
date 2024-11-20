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

func SetupTestMatchingInteractor(ctx context.Context, gw *testhelper.Gateway) (MatchingInteractor, *repository.UserMySQLRepository) {
	matchingRepo := repository.NewMatchingMySQLRepository(gw.MySQLClient) // ポインタを返すように &を追加
	userRepo := repository.NewUserMySQLRepository(gw.MySQLClient)
	ds := service.NewMatchingDomainService(userRepo, matchingRepo)
	return NewMatchingInteractor(matchingRepo, ds), &userRepo
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

func TestMatchingInteractor_Create(t *testing.T) {
	ctx := context.Background()
	gw, err := testhelper.Setup(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test: %v", err)
	}
	defer testhelper.Cleanup(ctx, gw)
	matchingInteractor, userRepo := SetupTestMatchingInteractor(ctx, gw)

	user1 := createTestUser(ctx, t, userRepo)
	user2 := createTestUser(ctx, t, userRepo)

	tests := []struct {
		name    string
		input   *input.CreateMatchingInput
		wantErr bool
	}{
		{
			name: "OK",
			input: &input.CreateMatchingInput{
				MeID:      user1.ID,
				PartnerID: user2.ID,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchingInteractor.Create(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.Matching == nil {
				t.Error("Create() got nil matching")
				return
			}
			if got.Matching.MeID != tt.input.MeID || got.Matching.PartnerID != tt.input.PartnerID {
				t.Errorf("Create() got = %v, want meID: %v, partnerID: %v", got.Matching, tt.input.MeID, tt.input.PartnerID)
			}
		})
	}
}

func TestMatchingInteractor_Get(t *testing.T) {
	ctx := context.Background()
	gw, err := testhelper.Setup(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test: %v", err)
	}
	defer testhelper.Cleanup(ctx, gw)
	matchingInteractor, userRepo := SetupTestMatchingInteractor(ctx, gw)

	user1 := createTestUser(ctx, t, userRepo)
	user2 := createTestUser(ctx, t, userRepo)

	created, err := matchingInteractor.Create(ctx, &input.CreateMatchingInput{
		MeID:      user1.ID,
		PartnerID: user2.ID,
	})
	if err != nil {
		t.Fatalf("Failed to create test matching: %v", err)
	}

	tests := []struct {
		name    string
		input   *input.GetMatchingInput
		want    *model.Matching
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &input.GetMatchingInput{ID: created.Matching.ID},
			want:    created.Matching,
			wantErr: false,
		},
		{
			name:    "NotFound",
			input:   &input.GetMatchingInput{ID: uuid.New()},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchingInteractor.Get(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			diff := cmp.Diff(
				got.Matching,
				tt.want,
				cmpopts.IgnoreFields(model.Matching{}, "ID", "CreatedAt", "UpdatedAt"),
			)
			if diff != "" {
				t.Errorf("Get() mismatching (-want +got):\n%s", diff)
			}
		})
	}
}

func TestMatchingInteractor_List(t *testing.T) {
	ctx := context.Background()
	gw, err := testhelper.Setup(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test: %v", err)
	}
	defer testhelper.Cleanup(ctx, gw)
	matchingInteractor, userRepo := SetupTestMatchingInteractor(ctx, gw)

	user1 := createTestUser(ctx, t, userRepo)
	user2 := createTestUser(ctx, t, userRepo)

	for i := 0; i < 3; i++ {
		partner := createTestUser(ctx, t, userRepo)
		_, err := matchingInteractor.Create(ctx, &input.CreateMatchingInput{
			MeID:      user1.ID,
			PartnerID: partner.ID,
		})
		if err != nil {
			t.Fatalf("Failed to create test matching: %v", err)
		}
	}

	tests := []struct {
		name    string
		input   *input.ListMatchingInput
		want    int
		wantErr bool
	}{
		{
			name: "OK_AllMatchinges",
			input: &input.ListMatchingInput{
				UserID: user1.ID,
				Limit:  10,
				Offset: 0,
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "OK_LimitedMatchinges",
			input: &input.ListMatchingInput{
				UserID: user1.ID,
				Limit:  2,
				Offset: 0,
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "OK_NoMatchinges",
			input: &input.ListMatchingInput{
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
			got, err := matchingInteractor.List(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got.Matchinges) != tt.want {
				t.Errorf("List() got = %v matchinges, want %v", len(got.Matchinges), tt.want)
			}
		})
	}
}

func TestMatchingInteractor_Update(t *testing.T) {
	ctx := context.Background()
	gw, err := testhelper.Setup(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test: %v", err)
	}
	defer testhelper.Cleanup(ctx, gw)
	matchingInteractor, userRepo := SetupTestMatchingInteractor(ctx, gw)

	user1 := createTestUser(ctx, t, userRepo)
	user2 := createTestUser(ctx, t, userRepo)

	created, err := matchingInteractor.Create(ctx, &input.CreateMatchingInput{
		MeID:      user1.ID,
		PartnerID: user2.ID,
	})
	if err != nil {
		t.Fatalf("Failed to create test matching: %v", err)
	}

	tests := []struct {
		name    string
		input   *input.UpdateMatchingInput
		want    model.MatchingStatus
		wantErr bool
	}{
		{
			name: "OK",
			input: &input.UpdateMatchingInput{
				ID:     created.Matching.ID,
				Status: model.MatchingStatusAccepted,
			},
			want:    model.MatchingStatusAccepted,
			wantErr: false,
		},
		{
			name: "NG_NotFound",
			input: &input.UpdateMatchingInput{
				ID:     uuid.New(),
				Status: model.MatchingStatusAccepted,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchingInteractor.Update(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.Matching.Status != tt.want {
				t.Errorf("Update() got status = %v, want %v", got.Matching.Status, tt.want)
			}
		})
	}
}

func TestMatchingInteractor_Delete(t *testing.T) {
	ctx := context.Background()
	gw, err := testhelper.Setup(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test: %v", err)
	}
	defer testhelper.Cleanup(ctx, gw)
	matchingInteractor, userRepo := SetupTestMatchingInteractor(ctx, gw)

	user1 := createTestUser(ctx, t, userRepo)
	user2 := createTestUser(ctx, t, userRepo)

	created, err := matchingInteractor.Create(ctx, &input.CreateMatchingInput{
		MeID:      user1.ID,
		PartnerID: user2.ID,
	})
	if err != nil {
		t.Fatalf("Failed to create test matching: %v", err)
	}

	tests := []struct {
		name    string
		input   *input.DeleteMatchingInput
		wantErr bool
	}{
		{
			name: "OK",
			input: &input.DeleteMatchingInput{
				ID: created.Matching.ID,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchingInteractor.Delete(ctx, tt.input)
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
