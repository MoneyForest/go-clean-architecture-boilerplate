package response

import (
	"encoding/json"
	"errors"
	"net/http"

	appError "github.com/MoneyForest/go-clean-boilerplate/pkg/error"
)

type ErrorResponse struct {
	Message string                 `json:"message"`
	Code    string                 `json:"code"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func WriteError(w http.ResponseWriter, err error) {
	var appErr *appError.AppError
	if errors.As(err, &appErr) {
		switch appErr.Code {
		case appError.InvalidArgument:
			WriteJSONError(w, http.StatusBadRequest, appErr)
		case appError.NotFound:
			WriteJSONError(w, http.StatusNotFound, appErr)
		case appError.AlreadyExists:
			WriteJSONError(w, http.StatusConflict, appErr)
		case appError.Unauthorized:
			WriteJSONError(w, http.StatusUnauthorized, appErr)
		case appError.PermissionDenied:
			WriteJSONError(w, http.StatusForbidden, appErr)
		case appError.PreconditionFailed:
			WriteJSONError(w, http.StatusPreconditionFailed, appErr)
		case appError.Critical:
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
	if appErr, ok := err.(*appError.AppError); ok {
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
		json.NewEncoder(w).Encode(data)
	}
}
