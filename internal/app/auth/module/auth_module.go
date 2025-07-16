package module

import (
	"payroll/internal/app/auth"
	"payroll/internal/app/auth/service"

	"github.com/google/wire"
)

var AuthModule = wire.NewSet(
	wire.Bind(new(auth.AuthService), new(*service.AuthService)),
	service.NewAuthService,
)
