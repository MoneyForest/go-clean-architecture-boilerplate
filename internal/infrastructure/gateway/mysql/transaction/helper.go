// internal/infrastructure/gateway/mysql/transaction/helper.go
package transaction

import (
	"context"
	"database/sql"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/gateway/mysql/sqlc"
)

// トランザクションコンテキストキーの型を定義
type txKey struct{}

// GetTx コンテキストからトランザクションを取得
func GetTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return nil
}

// GetQueries トランザクションの有無に応じて適切なQueriesインスタンスを返す
func GetQueries(ctx context.Context, defaultQueries *sqlc.Queries) *sqlc.Queries {
	if tx := GetTx(ctx); tx != nil {
		return defaultQueries.WithTx(tx)
	}
	return defaultQueries
}
