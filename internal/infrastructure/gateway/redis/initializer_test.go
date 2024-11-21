package redis

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestInitRedis(t *testing.T) {
	tests := []struct {
		name    string
		config  RedisConfig
		wantErr bool
	}{
		{
			name: "OK",
			config: RedisConfig{
				Environment: "test",
				Host:        "localhost",
				Port:        "6379",
				Password:    "password",
			},
			wantErr: false,
		},
		{
			name: "NG",
			config: RedisConfig{
				Environment: "test",
				Host:        "nonexistent",
				Port:        "6379",
				Password:    "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			client, err := InitRedis(ctx, tt.config)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if client == nil {
				t.Fatal("client is nil")
			}

			const testKey = "test_key"
			const testValue = "test_value"

			err = client.Set(ctx, testKey, testValue, 1*time.Minute).Err()
			if err != nil {
				t.Fatalf("failed to set value: %v", err)
			}

			got, err := client.Get(ctx, testKey).Result()
			if err != nil {
				t.Fatalf("failed to get value: %v", err)
			}

			if diff := cmp.Diff(testValue, got); diff != "" {
				t.Errorf("value mismatching (-want +got):\n%s", diff)
			}

			if err := client.Del(ctx, testKey).Err(); err != nil {
				t.Fatalf("failed to cleanup: %v", err)
			}
			if err := client.Close(); err != nil {
				t.Fatalf("failed to close client: %v", err)
			}
		})
	}
}
