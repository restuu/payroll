package service

import (
	"context"
	"fmt"
	"time"

	"payroll/internal/app/auth"
	"payroll/internal/app/auth/dto"
	"payroll/internal/app/common/apperror"
	"payroll/internal/app/common/role"
	"payroll/internal/infrastructure/config"
	"payroll/internal/infrastructure/database/postgres/repository"
	"payroll/pkg/jwt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	cfg                config.AuthConfig
	employeeRepository auth.AuthRepository
	jwt                jwt.JWT[*dto.JWTClaims]
}

func NewAuthService(
	cfg config.AuthConfig,
	employeeRepository auth.AuthRepository,
	jwt jwt.JWT[*dto.JWTClaims],
) *AuthService {
	return &AuthService{
		cfg:                cfg,
		employeeRepository: employeeRepository,
		jwt:                jwt,
	}
}

func (a *AuthService) RegisterNewEmployee(ctx context.Context, params dto.RegisterNewEmployeeParams) (err error) {
	existingEmployee, err := a.employeeRepository.FindEmployeeByUsername(ctx, params.Username)
	if err != nil {
		return fmt.Errorf("AuthService.RegisterNewEmployee failed to find employee by username %s: %w",
			params.Username, err)
	}

	if existingEmployee != nil {
		return fmt.Errorf("AuthService.RegisterNewEmployee found duplicate employee with username %s: %w",
			params.Username, apperror.ErrConflict)
	}

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

	var roleType repository.RoleType
	switch params.Role {
	case role.ADMIN:
		roleType = repository.RoleTypeADMIN
	default:
		roleType = repository.RoleTypeUSER
	}

	err = a.employeeRepository.InsertEmployeeRole(ctx, &repository.InsertEmployeeRoleParams{
		Username:  params.Username,
		RoleName:  roleType,
		CreatedBy: adminIDStr,
		UpdatedBy: adminIDStr,
	})
	if err != nil {
		return fmt.Errorf("AuthService.RegisterNewEmployee failed to insert employee role: %w", err)
	}

	return nil
}

func (a *AuthService) Login(ctx context.Context, req dto.LoginRequest) (res *dto.LoginResult, err error) {
	employee, err := a.employeeRepository.FindEmployeeByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Login failed to find employee by username %s: %w", req.Username, err)
	}

	if employee == nil {
		return nil, fmt.Errorf("AuthService.Login employee not found by username %s: %w", req.Username, apperror.ErrNotFound)
	}

	err = bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(req.Password))
	if err != nil {
		return nil, apperror.ErrUnauthorized.WithError(fmt.Errorf("AuthService.Login invalid password: %w", err))
	}

	expiresAt := time.Now().Add(a.cfg.JWTExpiry)

	claims := &dto.JWTClaims{
		JWTPayloads: dto.JWTPayloads{
			Username: employee.Username,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			ID:        uuid.New().String(),
		},
	}

	token, err := a.jwt.SignClaims(claims)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Login failed to sign token: %w", err)
	}

	res = &dto.LoginResult{
		AccessToken: token,
		ExpiresAt:   expiresAt.Format(time.RFC3339),
	}

	return res, nil
}
