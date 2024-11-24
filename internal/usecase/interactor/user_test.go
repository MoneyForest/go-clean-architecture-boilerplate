package interactor

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/gateway/mysql/repository"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/gateway/mysql/transaction"
	redisRepo "github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/gateway/redis/repository"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/gateway/sqs"
	sqsRepo "github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/gateway/sqs/repository"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/usecase/port"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/testhelper"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
)

func SetupTestUserInteractor(ctx context.Context, gw *testhelper.Gateway) UserInteractor {
	return NewUserInteractor(
		transaction.NewMySQLTransactionManager(gw.MySQLClient),
		repository.NewUserMySQLRepository(gw.MySQLClient),
		redisRepo.NewUserRedisRepository(gw.RedisClient),
		sqsRepo.NewSQSRepository(gw.SQSClient.Client, gw.SQSClient.QueueURLs[sqs.SQSKeySample]),
	)
}

func TestUserInteractor_Create(t *testing.T) {
	ctx := context.Background()
	gw, err := testhelper.Setup(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test: %v", err)
	}
	defer testhelper.Cleanup(ctx, gw)
	userInteractor := SetupTestUserInteractor(ctx, gw)

	tests := []struct {
		name    string
		input   *port.CreateUserInput
		want    *model.User
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &port.CreateUserInput{Email: "test@example.com"},
			want:    &model.User{Email: "test@example.com"},
			wantErr: false,
		},
		{
			name:    "NG",
			input:   &port.CreateUserInput{Email: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userInteractor.Create(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			user := got.User
			diff := cmp.Diff(
				user,
				tt.want,
				cmpopts.IgnoreFields(model.User{}, "ID", "CreatedAt", "UpdatedAt"),
			)
			if diff != "" {
				t.Errorf("Create() mismatching (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUserInteractor_Get(t *testing.T) {
	ctx := context.Background()
	gw, err := testhelper.Setup(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test: %v", err)
	}
	defer testhelper.Cleanup(ctx, gw)
	userInteractor := SetupTestUserInteractor(ctx, gw)

	createInput := &port.CreateUserInput{
		Email: "test@example.com",
	}
	created, err := userInteractor.Create(ctx, createInput)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name    string
		input   *port.GetUserInput
		want    *model.User
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &port.GetUserInput{ID: created.User.ID},
			want:    created.User,
			wantErr: false,
		},
		{
			name:    "NotFound",
			input:   &port.GetUserInput{ID: uuid.New()},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userInteractor.Get(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got == nil {
				if tt.want != nil {
					t.Errorf("Get() got = nil, want %v", tt.want)
				}
				return
			}
			diff := cmp.Diff(got.User, tt.want)
			if diff != "" {
				t.Errorf("Get() mismatching (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUserInteractor_List(t *testing.T) {
	ctx := context.Background()
	gw, err := testhelper.Setup(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test: %v", err)
	}
	defer testhelper.Cleanup(ctx, gw)
	userInteractor := SetupTestUserInteractor(ctx, gw)

	for i := range 3 {
		_, err := userInteractor.Create(ctx, &port.CreateUserInput{
			Email: fmt.Sprintf("test%d@example.com", i),
		})
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}

	tests := []struct {
		name    string
		input   *port.ListUserInput
		want    int
		wantErr bool
	}{
		{
			name: "OK_AllUsers",
			input: &port.ListUserInput{
				Limit:  10,
				Offset: 0,
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "OK_LimitedUsers",
			input: &port.ListUserInput{
				Limit:  2,
				Offset: 0,
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "OK_WithOffset",
			input: &port.ListUserInput{
				Limit:  10,
				Offset: 2,
			},
			want:    1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userInteractor.List(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if len(got.Users) != tt.want {
				t.Errorf("List() got = %v users, want %v", len(got.Users), tt.want)
			}
		})
	}
}

func TestUserInteractor_Update(t *testing.T) {
	ctx := context.Background()
	gw, err := testhelper.Setup(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test: %v", err)
	}
	defer testhelper.Cleanup(ctx, gw)
	userInteractor := SetupTestUserInteractor(ctx, gw)

	tests := []struct {
		name    string
		setup   func() (*model.User, error)
		input   *port.UpdateUserInput
		want    string
		wantErr bool
	}{
		{
			name: "OK",
			setup: func() (*model.User, error) {
				created, err := userInteractor.Create(ctx, &port.CreateUserInput{
					Email: "test@example.com",
				})
				return created.User, err
			},
			input:   &port.UpdateUserInput{Email: "updated@example.com"},
			want:    "updated@example.com",
			wantErr: false,
		},
		{
			name: "NG_InvalidEmail",
			setup: func() (*model.User, error) {
				created, err := userInteractor.Create(ctx, &port.CreateUserInput{
					Email: "test@example.com",
				})
				return created.User, err
			},
			input: &port.UpdateUserInput{
				Email: "",
			},
			wantErr: true,
		},
		{
			name: "NG_UserNotFound",
			setup: func() (*model.User, error) {
				return nil, nil
			},
			input: &port.UpdateUserInput{
				ID:    uuid.New(),
				Email: "notfound@example.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				user, err := tt.setup()
				if err != nil {
					t.Fatalf("Failed to setup test user: %v", err)
				}
				if user != nil {
					tt.input.ID = user.ID
				}
				if tt.wantErr {
					return
				}
				got, err := userInteractor.Update(ctx, tt.input)
				if err != nil {
					t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got.User.Email != tt.want {
					t.Errorf("Update() got = %v, want %v", got.User.Email, tt.want)
				}
			}
		})
	}
}

func TestUserInteractor_Delete(t *testing.T) {
	ctx := context.Background()
	gw, err := testhelper.Setup(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test: %v", err)
	}
	defer testhelper.Cleanup(ctx, gw)
	userInteractor := SetupTestUserInteractor(ctx, gw)

	tests := []struct {
		name    string
		setup   func() (*model.User, error)
		input   *port.DeleteUserInput
		wantErr bool
	}{
		{
			name: "OK",
			setup: func() (*model.User, error) {
				created, err := userInteractor.Create(ctx, &port.CreateUserInput{
					Email: "delete@example.com",
				})
				return created.User, err
			},
			input:   &port.DeleteUserInput{},
			wantErr: false,
		},
		{
			name:    "NG_UserNotFound",
			setup:   func() (*model.User, error) { return nil, nil },
			input:   &port.DeleteUserInput{ID: uuid.New()},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				user, err := tt.setup()
				if err != nil {
					t.Fatalf("Failed to setup test user: %v", err)
				}
				if user != nil {
					tt.input.ID = user.ID
				}
				if tt.wantErr {
					return
				}
				got, err := userInteractor.Delete(ctx, tt.input)
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
				result, err := userInteractor.Get(ctx, &port.GetUserInput{ID: tt.input.ID})
				if err != nil {
					t.Errorf("Unexpected error checking deleted user: %v", err)
					return
				}
				if result != nil && result.User != nil {
					t.Error("Delete() user still exists after deletion")
				}
			}
		})
	}
}

func TestUserInteractor_EnqueueUserDeletion(t *testing.T) {
	ctx := context.Background()
	gw, err := testhelper.Setup(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test: %v", err)
	}
	defer testhelper.Cleanup(ctx, gw)
	userInteractor := SetupTestUserInteractor(ctx, gw)

	testUser, err := userInteractor.Create(ctx, &port.CreateUserInput{
		Email: "process@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name    string
		input   *port.EnqueueUserDeletionInput
		wantErr bool
	}{
		{
			name: "OK_ExistingUser",
			input: &port.EnqueueUserDeletionInput{
				ID: testUser.User.ID,
			},
			wantErr: false,
		},
		{
			name: "OK_NonExistingUser",
			input: &port.EnqueueUserDeletionInput{
				ID: uuid.New(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userInteractor.EnqueueUserDeletion(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnqueueUserDeletion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.ID != tt.input.ID {
				t.Errorf("EnqueueUserDeletion() got = %v, want %v", got.ID, tt.input.ID)
			}
		})
	}
}

func TestUserInteractor_DequeueAndDeleteUser(t *testing.T) {
	ctx := context.Background()
	gw, err := testhelper.Setup(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test: %v", err)
	}
	defer testhelper.Cleanup(ctx, gw)
	userInteractor := SetupTestUserInteractor(ctx, gw)

	tests := []struct {
		name    string
		setup   func() error
		input   *port.DequeueAndDeleteUserInput
		want    int
		wantErr bool
	}{
		{
			name: "OK_ExistingUser",
			setup: func() error {
				testUser, err := userInteractor.Create(ctx, &port.CreateUserInput{
					Email: "process1@example.com",
				})
				if err != nil {
					return err
				}
				_, err = userInteractor.EnqueueUserDeletion(ctx, &port.EnqueueUserDeletionInput{
					ID: testUser.User.ID,
				})
				return err
			},
			input: &port.DequeueAndDeleteUserInput{
				BatchSize: 10,
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "OK_EmptyQueue",
			setup: func() error {
				return nil
			},
			input: &port.DequeueAndDeleteUserInput{
				BatchSize: 10,
			},
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 各テストケースのセットアップを実行
			if tt.setup != nil {
				if err := tt.setup(); err != nil {
					t.Fatalf("Failed to setup test: %v", err)
				}
			}

			got, err := userInteractor.DequeueAndDeleteUser(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DequeueAndDeleteUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.DeletedCount != tt.want {
				t.Errorf("DequeueAndDeleteUser() got = %v, want %v", got.DeletedCount, tt.want)
			}
		})
	}
}
