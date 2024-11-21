package marshaller

import (
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/controller/http/request"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/controller/http/response"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/usecase/port"
	"github.com/google/uuid"
)

// Input Marshalling
func ToCreateUserInput(req *request.CreateUserRequestBody) *port.CreateUserInput {
	return &port.CreateUserInput{
		Email: req.Email,
	}
}

func ToGetUserInput(req *request.GetUserParams) *port.GetUserInput {
	return &port.GetUserInput{
		ID: uuid.MustParse(req.ID),
	}
}

func ToListUsersInput(limit, offset int) *port.ListUserInput {
	return &port.ListUserInput{
		Limit:  limit,
		Offset: offset,
	}
}

func ToUpdateUserInput(req *request.UpdateUserRequestBody, id string) *port.UpdateUserInput {
	return &port.UpdateUserInput{
		ID:    uuid.MustParse(id),
		Email: req.Email,
	}
}

func ToDeleteUserInput(req *request.DeleteUserParams) *port.DeleteUserInput {
	return &port.DeleteUserInput{
		ID: uuid.MustParse(req.ID),
	}
}

// Output Marshalling
func ToUserResponse(user *model.User) response.UserResponse {
	return response.UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToCreateUserResponse(output *port.CreateUserOutput) response.CreateUserResponse {
	return response.CreateUserResponse(ToUserResponse(output.User))
}

func ToGetUserResponse(output *port.GetUserOutput) response.GetUserResponse {
	return response.GetUserResponse(ToUserResponse(output.User))
}

func ToListUsersResponse(output *port.ListUserOutput, limit, offset int) response.ListUsersResponse {
	users := make([]response.UserResponse, len(output.Users))
	for i, user := range output.Users {
		users[i] = ToUserResponse(user)
	}

	return response.ListUsersResponse{
		Users:     users,
		Total:     len(output.Users),
		Page:      (offset / limit) + 1,
		PageSize:  limit,
		TotalPage: (len(output.Users) + limit - 1) / limit,
	}
}

func ToUpdateUserResponse(output *port.UpdateUserOutput) response.UpdateUserResponse {
	return response.UpdateUserResponse(ToUserResponse(output.User))
}

func ToDeleteUserResponse(output *port.DeleteUserOutput) response.DeleteUserResponse {
	return response.DeleteUserResponse{
		ID: output.ID.String(),
	}
}
