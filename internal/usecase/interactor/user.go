package interactor

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/repository"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/transaction"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/usecase/port/input"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/usecase/port/output"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
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
	txManager transaction.Manager
	repo      repository.UserRepository
	cache     repository.UserCacheRepository
	msgQueue  repository.MessageQueueRepository
}

func NewUserInteractor(
	txManager transaction.Manager,
	repo repository.UserRepository,
	cache repository.UserCacheRepository,
	msgQueue repository.MessageQueueRepository,
) UserInteractor {
	return &userInteractor{
		repo:      repo,
		cache:     cache,
		msgQueue:  msgQueue,
		txManager: txManager,
	}
}

func (i *userInteractor) Create(ctx context.Context, input *input.CreateUserInput) (*output.CreateUserOutput, error) {
	user := model.NewUser(model.InputUserParams{ID: uuid.Nil(), Email: input.Email})
	if err := user.Validate(); err != nil {
		return nil, err
	}

	var createdUser *model.User
	err := i.txManager.DoInTx(ctx, func(ctx context.Context) error {
		var err error
		createdUser, err = i.repo.Create(ctx, user)
		return err
	})
	if err != nil {
		return nil, err
	}

	if err := i.cache.SetWithTTL(ctx, user, 3600*time.Second); err != nil {
		log.Printf("failed to set cache: %v\n", err)
	}

	return &output.CreateUserOutput{User: createdUser}, nil
}

func (i *userInteractor) Get(ctx context.Context, input *input.GetUserInput) (*output.GetUserOutput, error) {
	user, err := i.cache.Get(ctx, input.ID)
	if err == nil {
		return &output.GetUserOutput{User: user}, nil
	}
	user, err = i.repo.Get(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if err := i.cache.SetWithTTL(ctx, user, 3600*time.Second); err != nil {
		log.Printf("failed to set cache: %v\n", err)
	}

	return &output.GetUserOutput{User: user}, nil
}

func (i *userInteractor) List(ctx context.Context, input *input.ListUserInput) (*output.ListUserOutput, error) {
	users, err := i.repo.List(ctx, input.Limit, input.Offset)
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

	var updatedUser *model.User
	err := i.txManager.DoInTx(ctx, func(ctx context.Context) error {
		var err error
		updatedUser, err = i.repo.Update(ctx, user)
		return err
	})
	if err != nil {
		return nil, err
	}
	if err := i.cache.SetWithTTL(ctx, user, 3600*time.Second); err != nil {
		log.Printf("failed to set cache: %v\n", err)
	}
	return &output.UpdateUserOutput{User: updatedUser}, nil
}

func (i *userInteractor) Delete(ctx context.Context, input *input.DeleteUserInput) (*output.DeleteUserOutput, error) {
	var deletedID *uuid.UUID
	err := i.txManager.DoInTx(ctx, func(ctx context.Context) error {
		var err error
		deletedID, err = i.repo.Delete(ctx, input.ID)
		return err
	})
	if err != nil {
		return nil, err
	}
	if err := i.cache.Delete(ctx, input.ID); err != nil {
		log.Printf("failed to delete cache: %v\n", err)
	}
	return &output.DeleteUserOutput{ID: deletedID}, nil
}

func (i *userInteractor) ProcessMessage(ctx context.Context, input *input.ProcessMessageInput) (*output.ProcessMessageOutput, error) {
	userIDBytes, err := json.Marshal(input.ID)
	if err != nil {
		return nil, err
	}

	if err := i.msgQueue.SendMessage(ctx, &model.Message{
		Body: string(userIDBytes),
	}); err != nil {
		return nil, err
	}

	msgs, err := i.msgQueue.ReceiveMessage(ctx, &repository.ReceiveMessageOptions{
		MaxNumberOfMessages: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(msgs) == 0 {
		return nil, nil
	}

	msg := msgs[0]
	defer func() {
		if err := i.msgQueue.DeleteMessage(ctx, msg.ReceiptHandle); err != nil {
			log.Printf("Failed to delete message: %v", err)
		}
	}()

	var userID uuid.UUID
	if err := json.Unmarshal([]byte(msg.Body), &userID); err != nil {
		return nil, err
	}

	log.Printf("Dequeued user_id: %s\n", userID)
	return &output.ProcessMessageOutput{ID: &userID}, nil
}
