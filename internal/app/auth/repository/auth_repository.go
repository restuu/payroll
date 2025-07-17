package repository

import (
	"context"
	"database/sql"
	"errors"

	"payroll/internal/infrastructure/database/postgres/repository"
)

type AuthRepository struct {
	repository.Querier
}

func NewAuthRepository(querier repository.Querier) *AuthRepository {
	return &AuthRepository{Querier: querier}
}

func (a *AuthRepository) FindEmployeeByUsername(ctx context.Context, username string) (*repository.Employee, error) {
	res, err := a.Querier.FindEmployeeByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return res, nil
}
