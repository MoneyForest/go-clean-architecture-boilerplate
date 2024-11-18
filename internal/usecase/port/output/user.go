package output

import "github.com/MoneyForest/go-clean-boilerplate/internal/domain/model"

type GetUserOutput struct {
	User *model.User `json:"user"`
}

type ListUserOutput struct {
	Users []*model.User `json:"users"`
}
