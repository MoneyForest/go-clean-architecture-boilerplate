package request

import (
	"encoding/json"
	"net/http"
	"strconv"

	domainerr "github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/error"
	"github.com/go-chi/chi/v5"
)

type CreateUserRequestBody struct {
	Email string `json:"email"`
}

type GetUserParams struct {
	ID string `param:"id"`
}

type ListUsersQueryParams struct {
	Page     int    `query:"page"`
	PageSize int    `query:"pageSize"`
	SortBy   string `query:"sortBy"`
	Email    string `query:"email"`
}

type UpdateUserParams struct {
	ID string `param:"id"`
}

type UpdateUserRequestBody struct {
	Email string `json:"email"`
}

type DeleteUserParams struct {
	ID string `param:"id"`
}

// Request Decoding
func DecodeListUserRequest(r *http.Request) (int, int, error) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}
	return limit, offset, nil
}

func DecodeCreateUserRequest(r *http.Request) (*CreateUserRequestBody, error) {
	var req CreateUserRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, domainerr.NewDomainError(domainerr.InvalidArgument, "Invalid request body", err, nil)
	}
	return &req, nil
}

func DecodeGetUserRequest(r *http.Request) (*GetUserParams, error) {
	return &GetUserParams{
		ID: chi.URLParam(r, "id"),
	}, nil
}

func DecodeUpdateUserRequest(r *http.Request) (*UpdateUserRequestBody, error) {
	var req UpdateUserRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, domainerr.NewDomainError(domainerr.InvalidArgument, "Invalid request body", err, nil)
	}
	return &req, nil
}

func DecodeDeleteUserRequest(r *http.Request) (*DeleteUserParams, error) {
	return &DeleteUserParams{
		ID: chi.URLParam(r, "id"),
	}, nil
}
