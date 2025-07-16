package service

import (
	"context"
	"fmt"
	"time"

	"payroll/internal/app/attendance"
	"payroll/internal/infrastructure/database/postgres/repository"
)

type AttendanceService struct {
	attendanceRepository attendance.AttendanceRepository
}

func NewAttendanceService(attendanceRepository attendance.AttendanceRepository) *AttendanceService {
	return &AttendanceService{attendanceRepository: attendanceRepository}
}

func (a *AttendanceService) ClockIn(ctx context.Context, employeeID int64) (err error) {
	// TODO: check has clock in

	employeeIDStr := fmt.Sprintf("%d", employeeID)

	err = a.attendanceRepository.SaveEmployeeAttendance(ctx, &repository.SaveEmployeeAttendanceParams{
		EmployeeID: int32(employeeID),
		Timestamp:  time.Now(),
		Type:       repository.AttendanceTypeCLOCKIN,
		CreatedBy:  employeeIDStr,
		UpdatedBy:  employeeIDStr,
	})

	if err != nil {
		return fmt.Errorf("AttendanceService.ClockIn failed attendanceRepository.SaveEmployeeAttendance: %w", err)
	}

	return nil
}
