package repository

import "context"

type Transaction interface {
	DoInTx(ctx context.Context, fn func(ctx context.Context) error) error
}
