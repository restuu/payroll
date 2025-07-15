package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"payroll/internal/infrastructure/config"

	"github.com/go-chi/httplog/v3"
)

type requestIDCtxKey int

const RequestIDCtxKey requestIDCtxKey = 0

type ContextJSONHandler struct {
	jsonHandler slog.Handler
}

func (c *ContextJSONHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return c.jsonHandler.Enabled(ctx, level)
}

func (c *ContextJSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextJSONHandler{jsonHandler: c.jsonHandler.WithAttrs(attrs)}
}

func (c *ContextJSONHandler) WithGroup(name string) slog.Handler {
	return &ContextJSONHandler{jsonHandler: c.jsonHandler.WithGroup(name)}
}

func (c *ContextJSONHandler) Handle(ctx context.Context, r slog.Record) error {
	if requestID, ok := ctx.Value(RequestIDCtxKey).(string); ok {
		r.AddAttrs(slog.String("request_id", requestID))
	}

	return c.jsonHandler.Handle(ctx, r)
}

func SetDefaultLogger(cfg *config.Config) *slog.Logger {
	env := cfg.App.Env

	logFormat := httplog.SchemaOTEL

	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: logFormat.ReplaceAttr,
	})

	contextHandler := &ContextJSONHandler{
		jsonHandler: jsonHandler,
	}

	fmt.Println(".,.,.,.,.,.,.,.,.,")

	logger := slog.New(contextHandler).With(
		slog.String("app", cfg.App.Name),
		slog.String("version", "v1.0.0"),
		slog.String("env", env),
	)

	slog.SetDefault(logger)

	return logger
}
