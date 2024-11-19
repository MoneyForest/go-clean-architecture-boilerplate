package interactor

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/model"
	mysql "github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/mysql/repository"
	redis "github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/redis/repository"
	sqs "github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/gateway/sqs/repository"
	"github.com/MoneyForest/go-clean-boilerplate/internal/usecase/port/input"
	"github.com/MoneyForest/go-clean-boilerplate/internal/usecase/port/output"
	"github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"
)

type UserInteractor interface {
	Create(ctx context.Context, input *input.CreateUserInput) error
	Get(ctx context.Context, input *input.GetUserInput) (*output.GetUserOutput, error)
	List(ctx context.Context, input *input.ListUserInput) (*output.ListUserOutput, error)
	Update(ctx context.Context, input *input.UpdateUserInput) error
	Delete(ctx context.Context, input *input.DeleteUserInput) error
	ProcessMessage(ctx context.Context, input *input.ProcessMessageInput) error
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

func (i *userInteractor) Create(ctx context.Context, input *input.CreateUserInput) error {
	user := model.NewUser(model.InputUserParams{ID: uuid.Nil(), Email: input.Email})
	if err := user.Validate(); err != nil {
		return err
	}

	tx, err := i.mysql.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("failed to rollback transaction: %v\n", err)
		}
	}()

	if err := i.mysql.CreateTx(ctx, tx, user); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	if err := i.redis.SetWithTTL(ctx, user, 3600*time.Second); err != nil {
		log.Printf("failed to set cache: %v\n", err)
	}
	return nil
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

func (i *userInteractor) Update(ctx context.Context, input *input.UpdateUserInput) error {
	user := model.NewUser(model.InputUserParams{ID: input.ID, Email: input.Email})
	if err := user.Validate(); err != nil {
		return err
	}

	tx, err := i.mysql.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("failed to rollback transaction: %v\n", err)
		}
	}()

	if err := i.mysql.UpdateTx(ctx, tx, user); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	if err := i.redis.SetWithTTL(ctx, user, 3600*time.Second); err != nil {
		log.Printf("failed to set cache: %v\n", err)
	}
	return nil
}

func (i *userInteractor) Delete(ctx context.Context, input *input.DeleteUserInput) error {
	tx, err := i.mysql.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("failed to rollback transaction: %v\n", err)
		}
	}()

	if err := i.mysql.DeleteTx(ctx, tx, input.ID); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return i.redis.Delete(ctx, input.ID)
}

func (i *userInteractor) ProcessMessage(ctx context.Context, input *input.ProcessMessageInput) error {
	msgs, err := i.sqs.ReceiveMessage(ctx, input.QueueURL, &sqs.ReceiveMessageOptions{
		MaxNumberOfMessages: 1,
	})
	if err != nil {
		return err
	}
	if len(msgs) == 0 {
		return nil
	}

	msg := msgs[0]
	var userID uuid.UUID
	if err := json.Unmarshal([]byte(msg.Body), &userID); err != nil {
		return err
	}

	user, err := i.mysql.Get(ctx, userID)
	if err != nil {
		return err
	}
	if err := i.redis.SetWithTTL(ctx, user, 3600*time.Second); err != nil {
		log.Printf("failed to set cache: %v\n", err)
	}

	if err := i.sqs.DeleteMessage(ctx, input.QueueURL, msg.ReceiptHandle); err != nil {
		return err
	}

	log.Printf("Dequeued user_id: %s\n", userID)
	return nil
}
