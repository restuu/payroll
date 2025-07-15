package router

import (
	"log/slog"
	"net/http"

	"payroll/internal/infrastructure/config"
	"payroll/internal/presentation/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(cfg *config.Config, logger *slog.Logger) chi.Router {
	r := chi.NewRouter()

	r.Use(chimiddleware.RealIP)
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger(cfg, logger))
	r.Use(middleware.Cors())

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	return r
}
