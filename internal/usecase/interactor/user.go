package interactor

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/model"
	mysql "github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/mysql/repository"
	redis "github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/redis/repository"
	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/sqs/entity"
	sqs "github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/sqs/repository"
	"github.com/MoneyForest/go-clean-boilerplate/internal/usecase/port/input"
	"github.com/MoneyForest/go-clean-boilerplate/internal/usecase/port/output"
	"github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"
)

type UserInteractor interface {
	Create(ctx context.Context, input *input.CreateUserInput) (*output.CreateUserOutput, error)
	Get(ctx context.Context, input *input.GetUserInput) (*output.GetUserOutput, error)
	List(ctx context.Context, input *input.ListUserInput) (*output.ListUserOutput, error)
	Update(ctx context.Context, input *input.UpdateUserInput) (*output.UpdateUserOutput, error)
	Delete(ctx context.Context, input *input.DeleteUserInput) (*output.DeleteUserOutput, error)
	ProcessMessage(ctx context.Context, input *input.ProcessMessageInput) (*output.ProcessMessageOutput, error)
}

type userInteractor struct {
	mysql mysql.UserMySQLRepository
	redis redis.UserRedisRepository
	sqs   sqs.SQSRepository
}

func NewUserInteractor(mysql mysql.UserMySQLRepository, redis redis.UserRedisRepository, sqs sqs.SQSRepository) UserInteractor {
	return &userInteractor{
		mysql: mysql,
		redis: redis,
		sqs:   sqs,
	}
}

func (i *userInteractor) Create(ctx context.Context, input *input.CreateUserInput) (*output.CreateUserOutput, error) {
	user := model.NewUser(model.InputUserParams{ID: uuid.Nil(), Email: input.Email})
	if err := user.Validate(); err != nil {
		return nil, err
	}

	tx, err := i.mysql.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("failed to rollback transaction: %v\n", err)
		}
	}()

	createdUser, err := i.mysql.CreateTx(ctx, tx, user)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	if err := i.redis.SetWithTTL(ctx, user, 3600*time.Second); err != nil {
		log.Printf("failed to set cache: %v\n", err)
	}
	return &output.CreateUserOutput{User: createdUser}, nil
}

func (i *userInteractor) Get(ctx context.Context, input *input.GetUserInput) (*output.GetUserOutput, error) {
	user, err := i.redis.Get(ctx, input.ID)
	if err == nil {
		return &output.GetUserOutput{User: user}, nil
	}
	user, err = i.mysql.Get(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if err := i.redis.SetWithTTL(ctx, user, 3600*time.Second); err != nil {
		log.Printf("failed to set cache: %v\n", err)
	}

	return &output.GetUserOutput{User: user}, nil
}

func (i *userInteractor) List(ctx context.Context, input *input.ListUserInput) (*output.ListUserOutput, error) {
	users, err := i.mysql.List(ctx, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}
	return &output.ListUserOutput{Users: users}, nil
}

func (i *userInteractor) Update(ctx context.Context, input *input.UpdateUserInput) (*output.UpdateUserOutput, error) {
	user := model.NewUser(model.InputUserParams{ID: input.ID, Email: input.Email})
	if err := user.Validate(); err != nil {
		return nil, err
	}

	tx, err := i.mysql.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("failed to rollback transaction: %v\n", err)
		}
	}()

	updatedUser, err := i.mysql.UpdateTx(ctx, tx, user)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	if err := i.redis.SetWithTTL(ctx, user, 3600*time.Second); err != nil {
		log.Printf("failed to set cache: %v\n", err)
	}
	return &output.UpdateUserOutput{User: updatedUser}, nil
}

func (i *userInteractor) Delete(ctx context.Context, input *input.DeleteUserInput) (*output.DeleteUserOutput, error) {
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

	if err := i.redis.Delete(ctx, input.ID); err != nil {
		log.Printf("failed to delete cache: %v\n", err)
	}
	return &output.DeleteUserOutput{ID: deletedID}, nil
}

func (i *userInteractor) ProcessMessage(ctx context.Context, input *input.ProcessMessageInput) (*output.ProcessMessageOutput, error) {
	userIDBytes, err := json.Marshal(input.ID)
	if err != nil {
		return nil, err
	}

	if err := i.sqs.SendMessage(ctx, &entity.Message{
		Body: string(userIDBytes),
	}); err != nil {
		return nil, err
	}

	msgs, err := i.sqs.ReceiveMessage(ctx, &sqs.ReceiveMessageOptions{
		MaxNumberOfMessages: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(msgs) == 0 {
		return nil, nil
	}

	msg := msgs[0]
	var userID uuid.UUID
	if err := json.Unmarshal([]byte(msg.Body), &userID); err != nil {
		return nil, err
	}

	log.Printf("Dequeued user_id: %s\n", userID)
	return &output.ProcessMessageOutput{ID: &userID}, nil
}
