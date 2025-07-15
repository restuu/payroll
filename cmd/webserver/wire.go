//go:build wireinject
// +build wireinject

package main

import (
	"net/http"

	"payroll/internal/infrastructure/config"
	"payroll/internal/infrastructure/log"
	"payroll/internal/presentation/router"
	"payroll/internal/presentation/server"

	"github.com/go-chi/chi/v5"
	"github.com/google/wire"
)

func NewWebServer() (*WebServer, error) {
	wire.Build(
		config.LoadConfig,
		wire.FieldsOf(new(*config.Config), "Server"),
		log.SetDefaultLogger,
		router.NewRouter,
		wire.Bind(new(http.Handler), new(chi.Router)),
		server.NewServer,
		wire.Struct(new(WebServer), "*"),
	)

	return &WebServer{}, nil
}
