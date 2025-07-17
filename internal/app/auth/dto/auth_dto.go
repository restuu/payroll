package dto

import (
	"payroll/internal/app/common/role"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shopspring/decimal"
)

type RegisterNewEmployeeParams struct {
	Username string          `json:"username"`
	Password string          `json:"password"`
	Salary   decimal.Decimal `json:"salary"`
	Role     role.Role       `json:"role"`
	AdminID  int64           `json:"-"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResult struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   string `json:"expires_at"`
}

type JWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
