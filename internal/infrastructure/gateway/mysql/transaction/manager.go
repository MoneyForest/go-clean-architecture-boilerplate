package transaction

import (
	"context"
	"database/sql"
	"log"
)

type MySQLTransactionManager interface {
	DoInTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type mysqlTransactionManager struct {
	db *sql.DB
}

func NewMySQLTransactionManager(db *sql.DB) MySQLTransactionManager {
	return &mysqlTransactionManager{db: db}
}

func (t *mysqlTransactionManager) DoInTx(ctx context.Context, fn func(ctx context.Context) error) error {
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
