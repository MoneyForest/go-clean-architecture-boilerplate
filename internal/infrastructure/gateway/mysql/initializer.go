package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	Environment string
	Host        string
	Port        string
	User        string
	Password    string
	DBName      string
}

// InitDB はMySQLデータベース接続を初期化します
func InitDB(ctx context.Context, config DBConfig) (*sql.DB, error) {
	var dsn string

	switch config.Environment {
	case "local", "test":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
			config.User,
			config.Password,
			config.Host,
			config.Port,
			config.DBName,
		)
	default:
		return nil, fmt.Errorf("invalid environment: %s", config.Environment)
	}

	// データベース接続を確立
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mysql: %w", err)
	}

	// 接続テスト
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Successfully connected to MySQL at %s:%s\n", config.Host, config.Port)

	// コネクションプールの設定
	// 最大接続数の設定
	db.SetMaxOpenConns(25)
	// アイドル接続の最大数
	db.SetMaxIdleConns(25)
	// 接続の最大生存時間
	db.SetConnMaxLifetime(5 * time.Minute)
	// アイドル接続の最大生存時間
	db.SetConnMaxIdleTime(5 * time.Minute)

	return db, nil
}
