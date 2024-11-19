package interactor

import (
	"context"
	"testing"
	"time"

	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/model"
	mysqlRepo "github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/mysql/repository"
	redisRepo "github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/redis/repository"
	sqsRepo "github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/sqs/repository"
	"github.com/MoneyForest/go-clean-boilerplate/internal/usecase/port/input"
	"github.com/MoneyForest/go-clean-boilerplate/pkg/testhelper"
	"github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func SetupTestUserInteractor(ctx context.Context, gw *testhelper.Gateway) UserInteractor {
	return NewUserInteractor(
		mysqlRepo.NewUserMySQLRepository(gw.MySQLClient),
		redisRepo.NewUserRedisRepository(gw.RedisClient),
		sqsRepo.NewSQSRepository(gw.SQSClient),
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
			name: "OK",
			input: &input.CreateUserInput{
				Email: "test@example.com",
			},
			want: &model.User{
				ID:        uuid.Nil(),
				Email:     "test@example.com",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			wantErr: false,
		},
		{
			name: "NG",
			input: &input.CreateUserInput{
				Email: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userInteractor.Create(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return // エラーが期待される場合はここで終了
			}

			// エラーケースの場合はユーザー比較をスキップ
			if tt.wantErr {
				return
			}

			user := got.User
			diff := cmp.Diff(
				user, tt.want,
				cmpopts.IgnoreFields(model.User{}, "ID", "CreatedAt", "UpdatedAt"),
			)
			if diff != "" {
				t.Errorf("Create() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
