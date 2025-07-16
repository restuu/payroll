package auth

import (
	"context"

	"payroll/internal/app/auth/dto"
	"payroll/internal/infrastructure/database/postgres/repository"
)

type EmployeeRepository interface {
	InsertEmployee(ctx context.Context, arg *repository.InsertEmployeeParams) error
}

type AuthService interface {
	RegisterNewEmployee(ctx context.Context, params dto.RegisterNewEmployeeParams) (err error)
}
