//go:build wireinject
// +build wireinject

package main

import (
	"net/http"

	"payroll/internal/app"
	"payroll/internal/app/attendance"
	attendancemodule "payroll/internal/app/attendance/module"
	"payroll/internal/app/auth"
	authmodule "payroll/internal/app/auth/module"
	"payroll/internal/infrastructure/config"
	"payroll/internal/infrastructure/database/postgres"
	"payroll/internal/infrastructure/database/postgres/repository"
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
		wire.FieldsOf(new(*config.Config), "Database"),
		postgres.Connect,

		wire.Bind(new(attendance.AttendanceRepository), new(repository.Querier)),
		attendancemodule.AttendanceModule,
		wire.Bind(new(auth.EmployeeRepository), new(repository.Querier)),
		authmodule.AuthModule,

		wire.Struct(new(app.Services), "*"),

		router.NewRouter,
		wire.Bind(new(http.Handler), new(chi.Router)),
		server.NewServer,
		wire.Struct(new(WebServer), "*"),
	)

	return &WebServer{}, nil
}
