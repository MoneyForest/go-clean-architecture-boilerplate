package output

import (
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
)

type CreateMatchingOutput struct {
	Matching *model.Matching `json:"matching"`
}

type GetMatchingOutput struct {
	Matching *model.Matching `json:"matching"`
}

type ListMatchingOutput struct {
	Matchinges []*model.Matching `json:"matchinges"`
}

type UpdateMatchingOutput struct {
	Matching *model.Matching `json:"matching"`
}

type DeleteMatchingOutput struct {
	ID *uuid.UUID `json:"id"`
}
