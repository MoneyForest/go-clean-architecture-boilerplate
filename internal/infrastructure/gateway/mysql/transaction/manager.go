// internal/infrastructure/gateway/mysql/transaction/manager.go
package transaction

import (
	"context"
	"database/sql"
	"fmt"
)

type mysqlTransactionManager struct {
	db *sql.DB
}

func NewMySQLTransactionManager(db *sql.DB) *mysqlTransactionManager {
	return &mysqlTransactionManager{db: db}
}

func (m *mysqlTransactionManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	ctxWithTx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(ctxWithTx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback transaction: %v, original error: %w", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
