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

func (a *AttendanceService) ClockIn(ctx context.Context, username string) (err error) {
	// TODO: check has clock in
	now := time.Now()
	fmt.Println(now.Format(time.RFC3339))

	err = a.attendanceRepository.SaveEmployeeAttendance(ctx, &repository.SaveEmployeeAttendanceParams{
		Username:  username,
		Timestamp: time.Now(),
		Type:      repository.AttendanceTypeCLOCKIN,
		CreatedBy: username,
		UpdatedBy: username,
	})

	if err != nil {
		return fmt.Errorf("AttendanceService.ClockIn failed attendanceRepository.SaveEmployeeAttendance: %w", err)
	}

	return nil
}

func (a *AttendanceService) ClockOut(ctx context.Context, username string) (err error) {
	// TODO: check has clock out
	now := time.Now()
	fmt.Println(now.Format(time.RFC3339))

	err = a.attendanceRepository.SaveEmployeeAttendance(ctx, &repository.SaveEmployeeAttendanceParams{
		Username:  username,
		Timestamp: time.Now(),
		Type:      repository.AttendanceTypeCLOCKOUT,
		CreatedBy: username,
		UpdatedBy: username,
	})

	if err != nil {
		return fmt.Errorf("AttendanceService.ClockOut failed attendanceRepository.SaveEmployeeAttendance: %w", err)
	}

	return nil
}
