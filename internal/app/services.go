package app

import (
	"payroll/internal/app/attendance"
	"payroll/internal/app/auth"
	"payroll/internal/app/payroll"
)

type Services struct {
	AttendanceService attendance.AttendanceService
	AuthService       auth.AuthService
	PayrollService    payroll.PayrollService
}
