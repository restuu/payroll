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
	"payroll/internal/app/common/message"
	payrollmodule "payroll/internal/app/payroll/module"
	"payroll/internal/infrastructure/config"
	"payroll/internal/infrastructure/database/postgres"
	"payroll/internal/infrastructure/database/postgres/repository"
	"payroll/internal/infrastructure/log"
	"payroll/internal/infrastructure/messagebus"
	"payroll/internal/presentation"
	"payroll/internal/presentation/middleware"
	"payroll/internal/presentation/router"
	"payroll/internal/presentation/server"
	"payroll/pkg/jwt"

	"github.com/go-chi/chi/v5"
	"github.com/google/wire"
)

func ProvideJWT(cfg config.AuthConfig) jwt.JWT[*dto.JWTClaims] {
	return jwt.NewJWT[*dto.JWTClaims](jwt.Config{
		SecretKey: cfg.JWTSecret,
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

		wire.FieldsOf(new(*config.Config), "Kafka"),
		messagebus.NewKafkaClient,

		message.NewMessagePublisher,

		wire.Bind(new(attendance.AttendanceRepository), new(repository.Querier)),

		payrollmodule.PayrollModule,
		attendancemodule.AttendanceModule,
		authmodule.AuthModule,

		wire.Struct(new(app.Services), "*"),

		middleware.WithJWTAuth,
		middleware.IsAdmin,
		wire.Struct(new(presentation.Middlewares), "*"),

		router.NewRouter,
		wire.Bind(new(http.Handler), new(chi.Router)),
		server.NewServer,
		wire.Struct(new(WebServer), "*"),
	)

	return &WebServer{}, nil
}
