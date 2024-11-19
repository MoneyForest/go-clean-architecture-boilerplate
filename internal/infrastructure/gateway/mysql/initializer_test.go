package mysql

import (
	"context"
	"testing"
	"time"
)

func TestInitDB(t *testing.T) {
	tests := []struct {
		name        string
		config      DBConfig
		wantErr     bool
		errContains string
	}{
		{
			name: "OK",
			config: DBConfig{
				Environment: "local",
				Host:        "localhost",
				Port:        "3306",
				User:        "root",
				Password:    "password",
				DBName:      "maindb",
			},
			wantErr: false,
		},
		{
			name: "NG",
			config: DBConfig{
				Environment: "invalid",
				Host:        "localhost",
				Port:        "3306",
				User:        "root",
				Password:    "password",
				DBName:      "maindb",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			db, err := InitDB(ctx, tt.config)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if db == nil {
				t.Fatal("db is nil")
			}

			if db != nil {
				err := db.Close()
				if err != nil {
					t.Fatalf("failed to close db: %v", err)
				}
			}
		})
	}
}
