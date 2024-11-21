package transaction

import "context"

type Manager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
