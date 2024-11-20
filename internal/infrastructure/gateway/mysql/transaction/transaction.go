package transaction

import (
	"context"
	"database/sql"
	"log"
)

type Transaction interface {
	DoInTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type SQLTransaction struct {
	db *sql.DB
}

func NewSQLTransaction(db *sql.DB) Transaction {
	return &SQLTransaction{db: db}
}

func (t *SQLTransaction) DoInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("failed to rollback transaction: %v\n", err)
		}
	}()

	if err := fn(context.WithValue(ctx, "tx", tx)); err != nil {
		return err
	}

	return tx.Commit()
}
