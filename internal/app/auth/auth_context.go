package auth

import (
	"context"

	"payroll/internal/app/auth/dto"
)

type authContextKey struct{}

var AuthContextKey = &authContextKey{}

func WithAuthContext(ctx context.Context, payload dto.JWTPayloads) context.Context {
	return context.WithValue(ctx, AuthContextKey, payload)
}

func GetAuthContext(ctx context.Context) dto.JWTPayloads {
	payload, _ := ctx.Value(AuthContextKey).(dto.JWTPayloads)

	return payload
}
