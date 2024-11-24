package response

import (
	"time"
)

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateUserResponse UserResponse

type GetUserResponse UserResponse

type ListUsersResponse struct {
	Users     []UserResponse `json:"users"`
	Total     int            `json:"total"`
	Page      int            `json:"page"`
	PageSize  int            `json:"pageSize"`
	TotalPage int            `json:"totalPage"`
}

type UpdateUserResponse UserResponse

type DeleteUserResponse struct {
	ID string `json:"id"`
}
