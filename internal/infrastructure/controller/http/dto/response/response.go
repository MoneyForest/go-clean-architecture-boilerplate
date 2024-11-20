package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MoneyForest/go-clean-boilerplate/pkg/apperror"
)

type ErrorResponse struct {
	Message string                 `json:"message"`
	Code    string                 `json:"code"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func WriteError(w http.ResponseWriter, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		switch appErr.Code {
		case apperror.InvalidArgument:
			WriteJSONError(w, http.StatusBadRequest, appErr)
		case apperror.NotFound:
			WriteJSONError(w, http.StatusNotFound, appErr)
		case apperror.AlreadyExists:
			WriteJSONError(w, http.StatusConflict, appErr)
		case apperror.Unauthorized:
			WriteJSONError(w, http.StatusUnauthorized, appErr)
		case apperror.PermissionDenied:
			WriteJSONError(w, http.StatusForbidden, appErr)
		case apperror.PreconditionFailed:
			WriteJSONError(w, http.StatusPreconditionFailed, appErr)
		case apperror.Critical:
			WriteJSONError(w, http.StatusInternalServerError, appErr)
		default:
			WriteJSONError(w, http.StatusInternalServerError, appErr)
		}
		return
	}
	WriteJSONError(w, http.StatusInternalServerError, err)
}

func WriteJSONError(w http.ResponseWriter, status int, err error) {
	var response ErrorResponse
	if appErr, ok := err.(*apperror.AppError); ok {
		response = ErrorResponse{
			Message: appErr.Message,
			Code:    string(appErr.Code),
			Details: appErr.Details,
		}
	} else {
		response = ErrorResponse{
			Message: err.Error(),
			Code:    http.StatusText(status),
		}
	}

	WriteJSON(w, status, response)
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		}
	}
}
