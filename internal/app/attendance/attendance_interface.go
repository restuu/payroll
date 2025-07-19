package attendance

import (
	"context"

	"payroll/internal/infrastructure/database/postgres/repository"
)

type AttendanceRepository interface {
	SaveEmployeeAttendance(ctx context.Context, arg *repository.SaveEmployeeAttendanceParams) error
}

type AttendanceService interface {
	ClockIn(ctx context.Context, username string) (err error)
	ClockOut(ctx context.Context, username string) (err error)
}
