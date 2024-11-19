package sample

import (
	"context"
	"log"

	"github.com/MoneyForest/go-clean-boilerplate/internal/dependency"
	"github.com/MoneyForest/go-clean-boilerplate/internal/usecase/port/input"
	"github.com/google/uuid"
)

func Run(ctx context.Context, dependency *dependency.Dependency, args []string) error {
	userInteractor := dependency.UserInteractor
	user, err := userInteractor.Get(ctx, &input.GetUserInput{
		ID: uuid.MustParse(args[0]),
	})
	if err != nil {
		return err
	}
	log.Println(user)
	return nil
}
