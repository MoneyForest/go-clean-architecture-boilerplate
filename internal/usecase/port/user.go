package port

import (
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
)

type CreateUserInput struct {
	Email string `json:"email"`
}

type CreateUserOutput struct {
	User *model.User `json:"user"`
}

type GetUserInput struct {
	ID uuid.UUID `json:"id"`
}

type GetUserOutput struct {
	User *model.User `json:"user"`
}

type ListUserInput struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ListUserOutput struct {
	Users []*model.User `json:"users"`
}

type UpdateUserInput struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

type UpdateUserOutput struct {
	User *model.User `json:"user"`
}

type DeleteUserInput struct {
	ID uuid.UUID `json:"id"`
}

type DeleteUserOutput struct {
	ID *uuid.UUID `json:"id"`
}

type EnqueueUserDeletionInput struct {
	ID uuid.UUID
}

type EnqueueUserDeletionOutput struct {
	ID uuid.UUID
}

type DequeueAndDeleteUserInput struct {
	BatchSize int64
}

type DequeueAndDeleteUserOutput struct {
	DeletedCount int
}
