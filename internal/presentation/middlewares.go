package presentation

import (
	"payroll/internal/presentation/middleware"
)

type Middlewares struct {
	WithJWTAuth middleware.JWTAuthMiddleware
	IsAdmin     middleware.IsAdminMiddleware
}
