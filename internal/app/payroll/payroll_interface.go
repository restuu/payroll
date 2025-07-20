package payroll

import (
	"context"

	"payroll/internal/app/payroll/dto"
)

type PayrollService interface {
	SubmitGeneratePayrollTask(ctx context.Context, req dto.GeneratePayrollTaskRequest) (*dto.GeneratePayrollTaskResult, error)
}
