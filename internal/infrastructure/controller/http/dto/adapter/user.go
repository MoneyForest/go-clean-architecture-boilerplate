package adapter

import (
	"github.com/google/uuid"

	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/controller/http/dto"
	"github.com/MoneyForest/go-clean-boilerplate/internal/usecase/port/input"
	"github.com/MoneyForest/go-clean-boilerplate/internal/usecase/port/output"
)

func ToCreateUserInput(req *dto.CreateUserRequestBody) *input.CreateUserInput {
	return &input.CreateUserInput{
		Email: req.Email,
	}
}

func ToUpdateUserInput(req *dto.UpdateUserRequestBody, id string) *input.UpdateUserInput {
	return &input.UpdateUserInput{
		ID:    uuid.MustParse(id),
		Email: req.Email,
	}
}

func ToUserResponse(user *model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:    user.ID.String(),
		Email: user.Email,
	}
}

func ToListUsersResponse(output *output.ListUserOutput) []dto.UserResponse {
	users := make([]dto.UserResponse, len(output.Users))
	for i, user := range output.Users {
		users[i] = ToUserResponse(user)
	}
	return users
}
