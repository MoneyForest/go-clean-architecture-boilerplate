package logger

import (
	"context"
	"log/slog"
	"os"
	"time"
)

type CustomHandler struct {
	slog.Handler
}

func NewLogger(level slog.Level) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   a.Key,
					Value: slog.StringValue(a.Value.Time().Format(time.RFC3339)),
				}
			}
			return a
		},
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	customHandler := &CustomHandler{
		Handler: handler,
	}

	return slog.New(customHandler)
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	if traceID := ctx.Value("trace_id"); traceID != nil {
		r.Add("trace_id", traceID)
	}
	r.Add("caller", r.PC)
	return h.Handler.Handle(ctx, r)
}
