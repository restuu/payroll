package http

import (
	"log/slog"
	"net/http"

	"payroll/internal/app/attendance"
	"payroll/internal/app/auth"
	"payroll/internal/app/common/api"
	"payroll/internal/infrastructure/log"
	"payroll/internal/presentation"

	"github.com/go-chi/chi/v5"
)

type attendanceHTTPHandler struct {
	attendanceService attendance.AttendanceService
}

func RegisterAttendanceRoutes(
	r chi.Router,
	middlewares *presentation.Middlewares,
	attendanceService attendance.AttendanceService,
) {
	h := attendanceHTTPHandler{attendanceService}

	r.Route("/v1/attendance", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middlewares.WithJWTAuth)

			r.Post("/clock-in", api.HandleHTTP(h.clockIn))
			r.Post("/clock-out", api.HandleHTTP(h.clockOut))
		})
	})
}

func (h *attendanceHTTPHandler) clockIn(r *http.Request) (bool, error) {
	ctx := r.Context()

	acct := auth.GetAuthContext(ctx)

	if err := h.attendanceService.ClockIn(ctx, acct.Username); err != nil {
		slog.ErrorContext(ctx, "attendanceHTTPHandler.clockIn failed", log.WithErrorAttr(err))
		return false, err
	}

	return true, nil
}

func (h *attendanceHTTPHandler) clockOut(r *http.Request) (bool, error) {
	ctx := r.Context()

	acct := auth.GetAuthContext(ctx)

	if err := h.attendanceService.ClockOut(ctx, acct.Username); err != nil {
		slog.ErrorContext(ctx, "attendanceHTTPHandler.clockOut failed", log.WithErrorAttr(err))
		return false, err
	}

	return true, nil
}
