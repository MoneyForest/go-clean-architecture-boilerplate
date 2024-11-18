package subscriber

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/MoneyForest/go-clean-boilerplate/internal/dependency"
)

type MessageHandler func(ctx context.Context, dependency *dependency.Dependency, args []string) error

func Run(f MessageHandler, args []string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dependency, err := dependency.Inject(ctx)
	if err != nil {
		return err
	}

	return f(ctx, dependency, args)
}
