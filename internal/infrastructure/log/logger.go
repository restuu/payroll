package log

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"payroll/internal/infrastructure/config"

	"github.com/go-chi/httplog/v3"
)

type requestIDCtxKey int

const RequestIDCtxKey requestIDCtxKey = 0

func SetRequestIDContext(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, RequestIDCtxKey, reqID)
}

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

	logFormat := httplog.SchemaECS

	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: logFormat.ReplaceAttr,
	})

	contextHandler := &ContextJSONHandler{
		jsonHandler: jsonHandler,
	}

	logger := slog.New(contextHandler).With(
		slog.String("app", cfg.App.Name),
		slog.String("version", "v1.0.0"),
		slog.String("env", env),
	)

	slog.SetDefault(logger)

	return logger
}

func WithErrorAttr(err error) slog.Attr {
	return slog.String("error", err.Error())
}

func WithMetadata(metadata any) slog.Attr {
	b, _ := json.Marshal(metadata)

	return slog.String("metadata", string(b))
}
