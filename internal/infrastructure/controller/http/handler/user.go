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

type UserHandler struct {
	UserInteractor interactor.UserInteractor
}

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

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, apperror.NewAppError(apperror.InvalidArgument, "Invalid request body", err, nil))
		return
	}

	input := adapter.ToCreateUserInput(&req)
	if err := h.UserInteractor.Create(r.Context(), input); err != nil {
		response.WriteError(w, err)
		return
	}

	resp := dto.CreateUserResponse{
		Email: req.Email,
	}
	response.WriteJSON(w, http.StatusCreated, resp)
}

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
	if err := h.UserInteractor.Update(r.Context(), input); err != nil {
		response.WriteError(w, err)
		return
	}

	resp := dto.UpdateUserResponse{
		ID:    params.ID,
		Email: reqBody.Email,
	}
	response.WriteJSON(w, http.StatusOK, resp)
}

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

	if err := h.UserInteractor.Delete(r.Context(), input); err != nil {
		response.WriteError(w, err)
		return
	}

	resp := &dto.DeleteUserResponse{
		Success: true,
		Message: "User deleted successfully",
	}
	response.WriteJSON(w, http.StatusOK, resp)
}
