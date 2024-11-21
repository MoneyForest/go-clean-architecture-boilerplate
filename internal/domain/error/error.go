package error

import (
	"encoding/json"
	"fmt"
)

type ErrorCode string

const (
	InvalidArgument    ErrorCode = "INVALID_ARGUMENT"
	NotFound           ErrorCode = "NOT_FOUND"
	AlreadyExists      ErrorCode = "ALREADY_EXISTS"
	Unauthorized       ErrorCode = "UNAUTHORIZED"
	PermissionDenied   ErrorCode = "PERMISSION_DENIED"
	PreconditionFailed ErrorCode = "PRECONDITION_FAILED"
	Critical           ErrorCode = "CRITICAL"
)

type DomainError struct {
	Message string                 `json:"message"`
	Code    ErrorCode              `json:"code"`
	Cause   error                  `json:"cause,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func NewDomainError(code ErrorCode, message string, cause error, details map[string]interface{}) *DomainError {
	err := &DomainError{
		Message: message,
		Code:    code,
		Cause:   cause,
		Details: details,
	}

	return err
}

func (e *DomainError) Error() string {
	msg := fmt.Sprintf("%s (code: %s)", e.Message, e.Code)

	if e.Details != nil {
		if details, err := json.Marshal(e.Details); err == nil {
			msg += fmt.Sprintf(" | Details: %s", string(details))
		}
	}

	if e.Cause != nil {
		msg += fmt.Sprintf(" | Caused by: %v", e.Cause)
	}

	return msg
}

func (e *DomainError) Unwrap() error {
	return e.Cause
}
