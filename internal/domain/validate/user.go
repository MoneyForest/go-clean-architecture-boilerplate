package validate

import (
	e "github.com/MoneyForest/go-clean-boilerplate/internal/domain/error"
	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/model"
)

func ValidateUser(user *model.User) error {
	if user.Email == "" {
		return e.ErrUserEmailIsRequired
	}
	return nil
}
