package service

import (
	"context"
	"fmt"

	"payroll/internal/app/auth"
	"payroll/internal/app/auth/dto"
	"payroll/internal/infrastructure/database/postgres/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	employeeRepository auth.EmployeeRepository
}

func NewAuthService(
	employeeRepository auth.EmployeeRepository,
) *AuthService {
	return &AuthService{
		employeeRepository: employeeRepository,
	}
}

func (a *AuthService) RegisterNewEmployee(ctx context.Context, params dto.RegisterNewEmployeeParams) (err error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("AuthService.RegisterNewEmployee failed to encrypt password: %w", err)
	}

	adminIDStr := fmt.Sprintf("%d", params.AdminID)

	err = a.employeeRepository.InsertEmployee(ctx, &repository.InsertEmployeeParams{
		Username:  params.Username,
		Password:  string(encryptedPassword),
		Salary:    params.Salary,
		CreatedBy: adminIDStr,
		UpdatedBy: adminIDStr,
	})
	if err != nil {
		return fmt.Errorf("AuthService.RegisterNewEmployee failed to insert employee: %w", err)
	}

	return nil
}
