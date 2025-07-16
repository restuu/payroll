package module

import (
	"payroll/internal/app/attendance"
	"payroll/internal/app/attendance/service"

	"github.com/google/wire"
)

var AttendanceModule = wire.NewSet(
	wire.Bind(new(attendance.AttendanceService), new(*service.AttendanceService)),
	service.NewAttendanceService,
)
