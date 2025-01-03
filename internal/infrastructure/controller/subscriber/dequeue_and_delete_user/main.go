package dequeue_and_delete_user

import (
	"context"
	"log"
	"time"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/dependency"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/usecase/port"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
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
			if _, err := dependency.UserInteractor.DequeueAndDeleteUser(ctx, &port.DequeueAndDeleteUserInput{
				BatchSize: 100,
			}); err != nil {
				log.Printf("Error processing message: %v", err)
				time.Sleep(5 * time.Second)
			}
		}
	}
}
