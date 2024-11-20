package transaction

import "context"

type Manager interface {
	DoInTx(ctx context.Context, fn func(ctx context.Context) error) error
}
