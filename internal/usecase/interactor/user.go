package interactor

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/repository"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/usecase/port"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/usecase/transaction"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
)

type UserInteractor struct {
	txManager transaction.Manager
	userRepo  repository.UserRepository
	userCache repository.UserCacheRepository
	msgQueue  repository.MessageQueueRepository
}

func NewUserInteractor(
	txManager transaction.Manager,
	userRepo repository.UserRepository,
	userCache repository.UserCacheRepository,
	msgQueue repository.MessageQueueRepository,
) UserInteractor {
	return UserInteractor{
		txManager: txManager,
		userRepo:  userRepo,
		userCache: userCache,
		msgQueue:  msgQueue,
	}
}

func (i UserInteractor) Create(ctx context.Context, input *port.CreateUserInput) (*port.CreateUserOutput, error) {
	user := model.NewUser(model.InputUserParams{ID: uuid.Nil(), Email: input.Email})
	if err := user.Validate(); err != nil {
		return nil, err
	}

	var createdUser *model.User
	err := i.txManager.Do(ctx, func(ctx context.Context) error {
		var err error
		createdUser, err = i.userRepo.Save(ctx, user)
		return err
	})
	if err != nil {
		return nil, err
	}
	if createdUser != nil {
		if err := i.userCache.Store(ctx, createdUser, 3600*time.Second); err != nil {
			log.Printf("failed to set cache: %v\n", err)
		}
	}

	return &port.CreateUserOutput{User: createdUser}, nil
}

func (i UserInteractor) Get(ctx context.Context, input *port.GetUserInput) (*port.GetUserOutput, error) {
	user, err := i.userCache.FindById(ctx, input.ID)
	if err == nil {
		return &port.GetUserOutput{User: user}, nil
	}
	user, err = i.userRepo.FindById(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		if err := i.userCache.Store(ctx, user, 3600*time.Second); err != nil {
			log.Printf("failed to set cache: %v\n", err)
		}
	}
	return &port.GetUserOutput{User: user}, nil
}

func (i UserInteractor) List(ctx context.Context, input *port.ListUserInput) (*port.ListUserOutput, error) {
	users, err := i.userRepo.FindAll(ctx, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}
	return &port.ListUserOutput{Users: users}, nil
}

func (i UserInteractor) Update(ctx context.Context, input *port.UpdateUserInput) (*port.UpdateUserOutput, error) {
	user := model.NewUser(model.InputUserParams{ID: input.ID, Email: input.Email})
	if err := user.Validate(); err != nil {
		return nil, err
	}

	var updatedUser *model.User
	err := i.txManager.Do(ctx, func(ctx context.Context) error {
		var err error
		updatedUser, err = i.userRepo.Save(ctx, user)
		return err
	})
	if err != nil {
		return nil, err
	}
	if updatedUser != nil {
		if err := i.userCache.Store(ctx, updatedUser, 3600*time.Second); err != nil {
			log.Printf("failed to set cache: %v\n", err)
		}
	}
	return &port.UpdateUserOutput{User: updatedUser}, nil
}

func (i UserInteractor) Delete(ctx context.Context, input *port.DeleteUserInput) (*port.DeleteUserOutput, error) {
	var deletedID *uuid.UUID
	err := i.txManager.Do(ctx, func(ctx context.Context) error {
		var err error
		deletedID, err = i.userRepo.Remove(ctx, input.ID)
		return err
	})
	if err != nil {
		return nil, err
	}
	if deletedID != nil {
		if err := i.userCache.Remove(ctx, *deletedID); err != nil {
			log.Printf("failed to delete cache: %v\n", err)
		}
	}
	return &port.DeleteUserOutput{ID: deletedID}, nil
}

func (i UserInteractor) ProcessMessage(ctx context.Context, input *port.ProcessMessageInput) (*port.ProcessMessageOutput, error) {
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
	return &port.ProcessMessageOutput{ID: &userID}, nil
}
