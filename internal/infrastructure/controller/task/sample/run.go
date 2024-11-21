package sample

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/dependency"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/usecase/port"
)

func Run(ctx context.Context, dependency *dependency.Dependency, args []string) error {
	userInteractor := dependency.UserInteractor
	user, err := userInteractor.Get(ctx, &port.GetUserInput{
		ID: uuid.MustParse(args[0]),
	})
	if err != nil {
		return err
	}
	log.Println(user)
	return nil
}
