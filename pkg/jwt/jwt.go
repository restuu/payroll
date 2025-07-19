package jwt

import (
	"fmt"
	"reflect"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RegisteredClaims = jwt.RegisteredClaims

func NewNumericDate(t time.Time) *jwt.NumericDate {
	return jwt.NewNumericDate(t)
}

type Config struct {
	SecretKey string
}

type JWT[T jwt.Claims] interface {
	SignClaims(claims T) (string, error)
	ParseToken(tokenString string) (T, error)
}

type _jwt[T jwt.Claims] struct {
	cfg Config
}

func NewJWT[T jwt.Claims](cfg Config) JWT[T] {
	return &_jwt[T]{cfg: cfg}
}

func (j *_jwt[T]) SignClaims(claims T) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.cfg.SecretKey))
	if err != nil {
		return "", fmt.Errorf("JWT.SignClaims failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (j *_jwt[T]) ParseToken(tokenString string) (T, error) {
	var claims T

	// Initialize claims using reflection to avoid panicking if T is an interface
	// and to ensure a concrete type is used for parsing.
	// This assumes T is a pointer to a struct that implements jwt.Claims.
	claimsType := reflect.TypeFor[T]()
	claims = reflect.New(claimsType.Elem()).Interface().(T)

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			return []byte(j.cfg.SecretKey), nil
		},
	)
	if err != nil {
		return claims, fmt.Errorf("JWT.ParseToken failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(T)
	if !ok {
		return claims, fmt.Errorf("JWT.ParseToken unknown claims: %w", err)
	}

	return claims, nil
}
