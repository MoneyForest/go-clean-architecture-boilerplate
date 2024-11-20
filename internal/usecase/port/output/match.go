package output

import (
	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"
)

type CreateMatchOutput struct {
	Match *model.Match `json:"match"`
}

type GetMatchOutput struct {
	Match *model.Match `json:"match"`
}

type ListMatchOutput struct {
	Matches []*model.Match `json:"matches"`
}

type UpdateMatchOutput struct {
	Match *model.Match `json:"match"`
}

type DeleteMatchOutput struct {
	ID *uuid.UUID `json:"id"`
}
