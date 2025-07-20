package module

import (
	"payroll/internal/app/payroll"
	"payroll/internal/app/payroll/service"

	"github.com/google/wire"
)

var PayrollModule = wire.NewSet(
	wire.Bind(new(payroll.PayrollService), new(*service.PayrollService)),
	service.NewPayrollService,
)
