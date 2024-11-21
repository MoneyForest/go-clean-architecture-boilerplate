package output

import (
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
)

type CreateUserOutput struct {
	User *model.User `json:"user"`
}

type GetUserOutput struct {
	User *model.User `json:"user"`
}

type ListUserOutput struct {
	Users []*model.User `json:"users"`
}

type UpdateUserOutput struct {
	User *model.User `json:"user"`
}

type DeleteUserOutput struct {
	ID *uuid.UUID `json:"id"`
}

type ProcessMessageOutput struct {
	ID *uuid.UUID `json:"id"`
}
