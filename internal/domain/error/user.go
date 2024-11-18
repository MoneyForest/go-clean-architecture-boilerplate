package error

import "errors"

var (
	ErrUserEmailIsRequired = errors.New("user email is required")
)
