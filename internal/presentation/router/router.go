package router

import (
	"log/slog"
	"net/http"

	"payroll/internal/app"
	attendancehttp "payroll/internal/app/attendance/delivery/http"
	authhttp "payroll/internal/app/auth/delivery/http"
	payrollhttp "payroll/internal/app/payroll/delivery/http"
	"payroll/internal/infrastructure/config"
	"payroll/internal/presentation"
	"payroll/internal/presentation/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	cfg *config.Config,
	logger *slog.Logger,
	middlewares *presentation.Middlewares,
	services *app.Services,
) chi.Router {
	r := chi.NewRouter()

	r.Use(chimiddleware.RealIP)
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger(cfg, logger))
	r.Use(chimiddleware.StripSlashes)
	r.Use(middleware.Cors())

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Route("/api", func(r chi.Router) {
		authhttp.RegisterAuthRoutes(r, services.AuthService)
		attendancehttp.RegisterAttendanceRoutes(r, middlewares, services.AttendanceService)
		payrollhttp.RegisterPayrollRoutes(r, middlewares, services.PayrollService)
	})

	return r
}
