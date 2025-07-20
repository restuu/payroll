package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"payroll/internal/app/payroll"
	"payroll/internal/app/payroll/dto"

	"github.com/twmb/franz-go/pkg/kgo"
)

func HandlePayrollGenerate(record *kgo.Record, service payroll.PayrollService) error {
	var data dto.GeneratePayrollTaskRequest

	// unmarshal json value
	if err := json.Unmarshal(record.Value, &data); err != nil {
		return fmt.Errorf("HandlePayrollGenerate failed unmarshal json: %w", err)
	}

	ctx := record.Context

	slog.InfoContext(ctx, "HandlePayrollGenerate", "data", data)
	// TODO: handle create payroll docs in main business service and update task status

	return nil
}
