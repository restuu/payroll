package http

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"payroll/internal/app/common/api"
	"payroll/internal/app/payroll"
	"payroll/internal/app/payroll/dto"
	"payroll/internal/infrastructure/log"
	"payroll/internal/presentation"
	"payroll/internal/presentation/middleware"

	"github.com/go-chi/chi/v5"
)

type payrollHTTPHandler struct {
	payrollService payroll.PayrollService
}

func RegisterPayrollRoutes(
	r chi.Router,
	middlewares *presentation.Middlewares,
	payrollService payroll.PayrollService,
) {
	h := payrollHTTPHandler{payrollService: payrollService}

	r.Route("/v1/payroll", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(
				middlewares.WithJWTAuth,
				middlewares.IsAdmin,
			)

			r.Post("/generate", api.HandleHTTP(h.generate))
		})
	})
}

func (h *payrollHTTPHandler) generate(r *http.Request) (*dto.GeneratePayrollTaskResult, error) {
	ctx := r.Context()

	var req dto.GeneratePayrollTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("payrollHTTPHandler.generate failed to decode request body: %w", err)
	}

	if req.RequestID == "" {
		req.RequestID = middleware.GetRequestIDFromContext(ctx)
	}

	res, err := h.payrollService.SubmitGeneratePayrollTask(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, "payrollHTTPHandler.generate failed", log.WithErrorAttr(err))
		return nil, err
	}

	return res, nil
}
