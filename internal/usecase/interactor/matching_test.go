package interactor

import (
	"context"
	"testing"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/service"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/gateway/mysql/repository"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/gateway/mysql/transaction"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/usecase/port"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/testhelper"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
)

func SetupTestMatchingInteractor(ctx context.Context, gw *testhelper.Gateway) (MatchingInteractor, *repository.UserMySQLRepository) {
	txManager := transaction.NewMySQLTransactionManager(gw.MySQLClient)
	matchingRepo := repository.NewMatchingMySQLRepository(gw.MySQLClient)
	userRepo := repository.NewUserMySQLRepository(gw.MySQLClient)
	ds := &service.MatchingDomainService{}
	return NewMatchingInteractor(txManager, matchingRepo, userRepo, ds), userRepo
}

func createTestUser(ctx context.Context, t *testing.T, userRepo *repository.UserMySQLRepository) *model.User {
	id := uuid.New()
	user := model.NewUser(model.InputUserParams{
		ID:    id,
		Email: id.String() + "@example.com",
	})
	createdUser, err := userRepo.Save(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
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
		input   *port.CreateMatchingInput
		wantErr bool
	}{
		{
			name: "OK",
			input: &port.CreateMatchingInput{
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
