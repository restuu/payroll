//go:build wireinject
// +build wireinject

package main

import (
	"net/http"

	"payroll/internal/app"
	"payroll/internal/app/attendance"
	attendancemodule "payroll/internal/app/attendance/module"
	"payroll/internal/app/auth/dto"
	authmodule "payroll/internal/app/auth/module"
	"payroll/internal/infrastructure/config"
	"payroll/internal/infrastructure/database/postgres"
	"payroll/internal/infrastructure/database/postgres/repository"
	"payroll/internal/infrastructure/log"
	"payroll/internal/presentation"
	"payroll/internal/presentation/middleware"
	"payroll/internal/presentation/router"
	"payroll/internal/presentation/server"
	"payroll/pkg/jwt"

	"github.com/go-chi/chi/v5"
	"github.com/google/wire"
)

func ProvideJWT(cfg *config.Config) jwt.JWT[*dto.JWTClaims] {
	return jwt.NewJWT[*dto.JWTClaims](jwt.Config{
		SecretKey: cfg.Auth.JWTSecret,
	})
}

func NewWebServer() (*WebServer, error) {
	wire.Build(
		config.LoadConfig,
		wire.FieldsOf(new(*config.Config), "Server"),
		log.SetDefaultLogger,
		wire.FieldsOf(new(*config.Config), "Database"),
		postgres.Connect,
		wire.FieldsOf(new(*config.Config), "Auth"),

		ProvideJWT,

		middleware.WithJWTAuth,
		wire.Struct(new(presentation.Middlewares), "*"),

		wire.Bind(new(attendance.AttendanceRepository), new(repository.Querier)),

		attendancemodule.AttendanceModule,
		authmodule.AuthModule,

		wire.Struct(new(app.Services), "*"),

		router.NewRouter,
		wire.Bind(new(http.Handler), new(chi.Router)),
		server.NewServer,
		wire.Struct(new(WebServer), "*"),
	)

	return &WebServer{}, nil
}
