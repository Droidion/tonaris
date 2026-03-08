package logging

import (
	"context"
	"log/slog"
)

type loggerContextKey struct{}

func FromContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return slog.Default()
	}

	logger, ok := ctx.Value(loggerContextKey{}).(*slog.Logger)
	if !ok || logger == nil {
		return slog.Default()
	}

	return logger
}

func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if logger == nil {
		return ctx
	}

	return context.WithValue(ctx, loggerContextKey{}, logger)
}
