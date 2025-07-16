package app

import (
	"payroll/internal/app/attendance"
	"payroll/internal/app/auth"
)

type Services struct {
	AttendanceService attendance.AttendanceService
	AuthService       auth.AuthService
}
