package interactor

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/mysql/repository"
	redisRepo "github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/redis/repository"
	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/sqs"
	sqsRepo "github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/sqs/repository"
	"github.com/MoneyForest/go-clean-boilerplate/internal/usecase/port/input"
	"github.com/MoneyForest/go-clean-boilerplate/pkg/testhelper"
	"github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"
)

func SetupTestUserInteractor(ctx context.Context, gw *testhelper.Gateway) UserInteractor {
	return NewUserInteractor(
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
		input   *input.CreateUserInput
		want    *model.User
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &input.CreateUserInput{Email: "test@example.com"},
			want:    &model.User{Email: "test@example.com"},
			wantErr: false,
		},
		{
			name:    "NG",
			input:   &input.CreateUserInput{Email: ""},
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

	createInput := &input.CreateUserInput{
		Email: "test@example.com",
	}
	created, err := userInteractor.Create(ctx, createInput)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name    string
		input   *input.GetUserInput
		want    *model.User
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &input.GetUserInput{ID: created.User.ID},
			want:    created.User,
			wantErr: false,
		},
		{
			name:    "NotFound",
			input:   &input.GetUserInput{ID: uuid.New()},
			want:    nil,
			wantErr: true,
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
		_, err := userInteractor.Create(ctx, &input.CreateUserInput{
			Email: fmt.Sprintf("test%d@example.com", i),
		})
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}

	tests := []struct {
		name    string
		input   *input.ListUserInput
		want    int
		wantErr bool
	}{
		{
			name: "OK_AllUsers",
			input: &input.ListUserInput{
				Limit:  10,
				Offset: 0,
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "OK_LimitedUsers",
			input: &input.ListUserInput{
				Limit:  2,
				Offset: 0,
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "OK_WithOffset",
			input: &input.ListUserInput{
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
		input   *input.UpdateUserInput
		want    string
		wantErr bool
	}{
		{
			name: "OK",
			setup: func() (*model.User, error) {
				created, err := userInteractor.Create(ctx, &input.CreateUserInput{
					Email: "test@example.com",
				})
				return created.User, err
			},
			input:   &input.UpdateUserInput{Email: "updated@example.com"},
			want:    "updated@example.com",
			wantErr: false,
		},
		{
			name: "NG_InvalidEmail",
			setup: func() (*model.User, error) {
				created, err := userInteractor.Create(ctx, &input.CreateUserInput{
					Email: "test@example.com",
				})
				return created.User, err
			},
			input: &input.UpdateUserInput{
				Email: "",
			},
			wantErr: true,
		},
		{
			name: "NG_UserNotFound",
			setup: func() (*model.User, error) {
				return nil, nil
			},
			input: &input.UpdateUserInput{
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
		input   *input.DeleteUserInput
		wantErr bool
	}{
		{
			name: "OK",
			setup: func() (*model.User, error) {
				created, err := userInteractor.Create(ctx, &input.CreateUserInput{
					Email: "delete@example.com",
				})
				return created.User, err
			},
			input:   &input.DeleteUserInput{},
			wantErr: false,
		},
		{
			name:    "NG_UserNotFound",
			setup:   func() (*model.User, error) { return nil, nil },
			input:   &input.DeleteUserInput{ID: uuid.New()},
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
				_, err = userInteractor.Get(ctx, &input.GetUserInput{ID: tt.input.ID})
				if err == nil {
					t.Error("Delete() user still exists after deletion")
				}
			}
		})
	}
}

func TestUserInteractor_ProcessMessage(t *testing.T) {
	ctx := context.Background()
	gw, err := testhelper.Setup(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test: %v", err)
	}
	defer testhelper.Cleanup(ctx, gw)
	userInteractor := SetupTestUserInteractor(ctx, gw)

	tests := []struct {
		name    string
		input   *input.ProcessMessageInput
		wantErr bool
	}{
		{
			name: "OK",
			input: &input.ProcessMessageInput{
				ID: uuid.New(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userInteractor.ProcessMessage(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got == nil || got.ID == nil {
				t.Error("ProcessMessage() got nil response")
				return
			}
			if *got.ID != tt.input.ID {
				t.Errorf("ProcessMessage() got = %v, want %v", *got.ID, tt.input.ID)
			}
		})
	}
}
