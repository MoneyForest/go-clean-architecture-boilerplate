package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/controller/http/marshaller"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/controller/http/request"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/controller/http/response"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/usecase/interactor"
)

// @title			User Handler
// @description	Handles HTTP requests for user operations
type UserHandler struct {
	UserInteractor interactor.UserInteractor
}

// @Summary	Create a new user
// @Tags		users
// @Accept		json
// @Produce	json
// @Param		body	body		request.CreateUserRequestBody	true	"User data"
// @Success	201		{object}	response.CreateUserResponse
// @Failure	400		{object}	apperror.AppError
// @Failure	500		{object}	apperror.AppError
// @Router		/users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	reqBody, err := request.DecodeCreateUserRequest(r)
	if err != nil {
		response.WriteError(w, err)
		return
	}
	output, err := h.UserInteractor.Create(
		r.Context(),
		marshaller.ToCreateUserInput(reqBody),
	)
	if err != nil {
		response.WriteError(w, err)
		return
	}
	response.WriteJSON(
		w,
		http.StatusCreated,
		marshaller.ToCreateUserResponse(output),
	)
}

// @Summary	Get user by ID
// @Tags		users
// @Accept		json
// @Produce	json
// @Param		id	path		string	true	"User ID"	format(uuid)
// @Success	200	{object}	response.GetUserResponse
// @Failure	400	{object}	apperror.AppError
// @Failure	404	{object}	apperror.AppError
// @Failure	500	{object}	apperror.AppError
// @Router		/users/{id} [get]
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	params, err := request.DecodeGetUserRequest(r)
	if err != nil {
		response.WriteError(w, err)
		return
	}
	output, err := h.UserInteractor.Get(
		r.Context(),
		marshaller.ToGetUserInput(params),
	)
	if err != nil {
		response.WriteError(w, err)
		return
	}
	response.WriteJSON(
		w,
		http.StatusOK,
		marshaller.ToGetUserResponse(output),
	)
}

// @Summary	List all users
// @Tags		users
// @Accept		json
// @Produce	json
// @Param		limit	query		int	false	"Items per page"	default(10)
// @Param		offset	query		int	false	"Skip items"		default(0)
// @Success	200		{object}	response.ListUsersResponse
// @Failure	400		{object}	apperror.AppError
// @Failure	500		{object}	apperror.AppError
// @Router		/users [get]
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := request.DecodeListUserRequest(r)
	if err != nil {
		response.WriteError(w, err)
		return
	}
	output, err := h.UserInteractor.List(
		r.Context(),
		marshaller.ToListUsersInput(limit, offset),
	)
	if err != nil {
		response.WriteError(w, err)
		return
	}
	response.WriteJSON(
		w,
		http.StatusOK,
		marshaller.ToListUsersResponse(output, limit, offset),
	)
}

// @Summary	Update user by ID
// @Tags		users
// @Accept		json
// @Produce	json
// @Param		id		path		string							true	"User ID"	format(uuid)
// @Param		body	body		request.UpdateUserRequestBody	true	"User data"
// @Success	200		{object}	response.UpdateUserResponse
// @Failure	400		{object}	apperror.AppError
// @Failure	404		{object}	apperror.AppError
// @Failure	500		{object}	apperror.AppError
// @Router		/users/{id} [put]
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	params := &request.UpdateUserParams{
		ID: chi.URLParam(r, "id"),
	}
	reqBody, err := request.DecodeUpdateUserRequest(r)
	if err != nil {
		response.WriteError(w, err)
		return
	}
	output, err := h.UserInteractor.Update(
		r.Context(),
		marshaller.ToUpdateUserInput(reqBody, params.ID),
	)
	if err != nil {
		response.WriteError(w, err)
		return
	}
	response.WriteJSON(
		w,
		http.StatusOK,
		marshaller.ToUpdateUserResponse(output),
	)
}

// @Summary	Delete user by ID
// @Tags		users
// @Accept		json
// @Produce	json
// @Param		id	path		string	true	"User ID"	format(uuid)
// @Success	200	{object}	response.DeleteUserResponse
// @Failure	400	{object}	apperror.AppError
// @Failure	404	{object}	apperror.AppError
// @Failure	500	{object}	apperror.AppError
// @Router		/users/{id} [delete]
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	params, err := request.DecodeDeleteUserRequest(r)
	if err != nil {
		response.WriteError(w, err)
		return
	}
	output, err := h.UserInteractor.Delete(
		r.Context(),
		marshaller.ToDeleteUserInput(params),
	)
	if err != nil {
		response.WriteError(w, err)
		return
	}
	response.WriteJSON(
		w,
		http.StatusOK,
		marshaller.ToDeleteUserResponse(output),
	)
}
