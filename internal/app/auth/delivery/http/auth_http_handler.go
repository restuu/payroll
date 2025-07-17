package http

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"payroll/internal/app/auth"
	"payroll/internal/app/auth/dto"
	"payroll/internal/app/common/api"
	"payroll/internal/infrastructure/log"

	"github.com/go-chi/chi/v5"
)

type authHTTPHandler struct {
	authService auth.AuthService
}

func RegisterAuthRoutes(r chi.Router, authService auth.AuthService) {
	h := &authHTTPHandler{authService: authService}

	r.Route("/v1", func(r chi.Router) {
		r.Post("/register", api.HandleHTTP(h.registerEmployee))
		r.Post("/login", api.HandleHTTP(h.login))
	})
}

func (h *authHTTPHandler) registerEmployee(r *http.Request) (any, error) {
	var data dto.RegisterNewEmployeeParams
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	err = h.authService.RegisterNewEmployee(r.Context(), data)
	if err != nil {
		slog.ErrorContext(r.Context(), "authHTTPHandler.registerEmployee failed", log.WithErrorAttr(err))
		return nil, err
	}

	return nil, nil
}

func (h *authHTTPHandler) login(r *http.Request) (*dto.LoginResult, error) {
	var data dto.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("authHTTPHandler.login failed to decode request body: %w", err)
	}

	res, err := h.authService.Login(r.Context(), data)
	if err != nil {
		slog.ErrorContext(r.Context(), "authHTTPHandler.login failed", log.WithErrorAttr(err))
		return nil, err
	}

	return res, nil
}
