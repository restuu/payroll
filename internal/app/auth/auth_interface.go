package auth

import (
	"context"

	"payroll/internal/app/auth/dto"
	"payroll/internal/infrastructure/database/postgres/repository"
)

type AuthRepository interface {
	InsertEmployee(ctx context.Context, arg *repository.InsertEmployeeParams) error
	FindEmployeeByUsername(ctx context.Context, username string) (*repository.Employee, error)
	InsertEmployeeRole(ctx context.Context, arg *repository.InsertEmployeeRoleParams) error
}

type AuthService interface {
	RegisterNewEmployee(ctx context.Context, params dto.RegisterNewEmployeeParams) (err error)
	Login(ctx context.Context, req dto.LoginRequest) (res *dto.LoginResult, err error)
}
