package handler

import (
	"encoding/json"
	"net/http"

	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/controller/http/dto"
	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/controller/http/dto/adapter"
	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/controller/http/dto/request"
	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/controller/http/dto/response"
	"github.com/MoneyForest/go-clean-boilerplate/internal/usecase/interactor"
	"github.com/MoneyForest/go-clean-boilerplate/internal/usecase/port/input"
	"github.com/MoneyForest/go-clean-boilerplate/pkg/apperror"
	"github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"
	"github.com/go-chi/chi/v5"
)

// @title			User Handler
// @description	Handles HTTP requests for user operations
type UserHandler struct {
	UserInteractor interactor.UserInteractor
}

// @Summary		List users
// @Description	Get a paginated list of users
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			limit	query		int	false	"Number of items per page"	default(10)
// @Param			offset	query		int	false	"Number of items to skip"	default(0)
// @Success		200		{object}	dto.ListUsersResponse
// @Failure		400		{object}	apperror.AppError
// @Failure		500		{object}	apperror.AppError
// @Router			/users [get]
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := request.ParseListParams(r)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	input := &input.ListUserInput{
		Limit:  limit,
		Offset: offset,
	}

	output, err := h.UserInteractor.List(r.Context(), input)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	resp := &dto.ListUsersResponse{
		Users:     adapter.ToListUsersResponse(output),
		Total:     len(output.Users),
		Page:      (offset / limit) + 1,
		PageSize:  limit,
		TotalPage: (len(output.Users) + limit - 1) / limit,
	}
	response.WriteJSON(w, http.StatusOK, resp)
}

// @Summary		Create user
// @Description	Create a new user
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			request	body		dto.CreateUserRequestBody	true	"User creation request"
// @Success		201		{object}	dto.CreateUserResponse
// @Failure		400		{object}	apperror.AppError
// @Failure		500		{object}	apperror.AppError
// @Router			/users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, apperror.NewAppError(apperror.InvalidArgument, "Invalid request body", err, nil))
		return
	}

	input := adapter.ToCreateUserInput(&req)
	_, err := h.UserInteractor.Create(r.Context(), input)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	resp := dto.CreateUserResponse{
		Email: req.Email,
	}
	response.WriteJSON(w, http.StatusCreated, resp)
}

// @Summary		Get user
// @Description	Get user by ID
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			id	path		string	true	"User ID"	format(uuid)
// @Success		200	{object}	dto.GetUserResponse
// @Failure		400	{object}	apperror.AppError
// @Failure		404	{object}	apperror.AppError
// @Failure		500	{object}	apperror.AppError
// @Router			/users/{id} [get]
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	params := &dto.GetUserParams{
		ID: chi.URLParam(r, "id"),
	}

	if !uuid.IsValidUUIDv7(params.ID) {
		response.WriteError(w, apperror.NewAppError(apperror.InvalidArgument, "Invalid user ID", nil, nil))
		return
	}

	input := &input.GetUserInput{
		ID: uuid.MustParse(params.ID),
	}

	output, err := h.UserInteractor.Get(r.Context(), input)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	resp := dto.GetUserResponse(adapter.ToUserResponse(output.User))
	response.WriteJSON(w, http.StatusOK, resp)
}

// @Summary		Update user
// @Description	Update user by ID
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			id		path		string						true	"User ID"	format(uuid)
// @Param			request	body		dto.UpdateUserRequestBody	true	"User update request"
// @Success		200		{object}	dto.UpdateUserResponse
// @Failure		400		{object}	apperror.AppError
// @Failure		404		{object}	apperror.AppError
// @Failure		500		{object}	apperror.AppError
// @Router			/users/{id} [put]
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	params := &dto.UpdateUserParams{
		ID: chi.URLParam(r, "id"),
	}

	if !uuid.IsValidUUIDv7(params.ID) {
		response.WriteError(w, apperror.NewAppError(apperror.InvalidArgument, "Invalid user ID", nil, nil))
		return
	}

	var reqBody dto.UpdateUserRequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		response.WriteError(w, apperror.NewAppError(apperror.InvalidArgument, "Invalid request body", err, nil))
		return
	}

	input := adapter.ToUpdateUserInput(&reqBody, params.ID)
	_, err := h.UserInteractor.Update(r.Context(), input)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	resp := dto.UpdateUserResponse{
		ID:    params.ID,
		Email: reqBody.Email,
	}
	response.WriteJSON(w, http.StatusOK, resp)
}

// @Summary		Delete user
// @Description	Delete user by ID
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			id	path		string	true	"User ID"	format(uuid)
// @Success		200	{object}	dto.DeleteUserResponse
// @Failure		400	{object}	apperror.AppError
// @Failure		404	{object}	apperror.AppError
// @Failure		500	{object}	apperror.AppError
// @Router			/users/{id} [delete]
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	params := &dto.DeleteUserParams{
		ID: chi.URLParam(r, "id"),
	}

	if !uuid.IsValidUUIDv7(params.ID) {
		response.WriteError(w, apperror.NewAppError(apperror.InvalidArgument, "Invalid user ID", nil, nil))
		return
	}

	input := &input.DeleteUserInput{
		ID: uuid.MustParse(params.ID),
	}

	_, err := h.UserInteractor.Delete(r.Context(), input)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	resp := &dto.DeleteUserResponse{
		Success: true,
		Message: "User deleted successfully",
	}
	response.WriteJSON(w, http.StatusOK, resp)
}
