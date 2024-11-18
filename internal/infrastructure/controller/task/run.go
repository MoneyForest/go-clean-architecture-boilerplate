package task

import (
	"context"
	"time"

	"github.com/MoneyForest/go-clean-boilerplate/internal/dependency"
)

func Run(f func(ctx context.Context, dependency *dependency.Dependency, args []string) error, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Second)
	defer cancel()

	dependency, err := dependency.Inject(ctx)
	if err != nil {
		return err
	}
	return f(ctx, dependency, args)
}
