package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type LoggerKey string

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		traceID := uuid.New().String()

		ctx := r.Context()
		logger := slog.With(
			"trace_id", traceID,
			"method", r.Method,
			"path", r.URL.Path,
		)
		ctx = context.WithValue(ctx, LoggerKey("logger"), logger)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

		logger.Info("request completed",
			"duration", time.Since(start),
		)
	})
}
