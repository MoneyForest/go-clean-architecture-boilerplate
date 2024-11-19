package sample

import (
	"context"
	"log"
	"time"

	"github.com/MoneyForest/go-clean-boilerplate/internal/dependency"
	"github.com/MoneyForest/go-clean-boilerplate/internal/usecase/port/input"
	"github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"
)

type Message struct {
	UserID uuid.UUID `json:"user_id"`
}

func Run(ctx context.Context, dependency *dependency.Dependency, args []string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if _, err := dependency.UserInteractor.ProcessMessage(ctx, &input.ProcessMessageInput{
				ID: uuid.New(),
			}); err != nil {
				log.Printf("Error processing message: %v", err)
				time.Sleep(5 * time.Second)
			}
		}
	}
}