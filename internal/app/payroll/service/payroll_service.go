package service

import (
	"context"

	"payroll/internal/app/common/message"
	"payroll/internal/app/payroll/dto"
	"payroll/internal/presentation/messaging/topics"

	"github.com/google/uuid"
)

type PayrollService struct {
	publisher message.MessagePublisher
}

func NewPayrollService(publisher message.MessagePublisher) *PayrollService {
	return &PayrollService{publisher: publisher}
}

func (p *PayrollService) SubmitGeneratePayrollTask(ctx context.Context, req dto.GeneratePayrollTaskRequest) (*dto.GeneratePayrollTaskResult, error) {
	// TODO: check running job
	// TODO: insert task request to db
	if req.RequestID == "" {
		req.RequestID = uuid.NewString()
	}

	err := p.publisher.PublishJSON(ctx, topics.PayrollGenerate, req)
	if err != nil {
		return nil, err
	}

	result := &dto.GeneratePayrollTaskResult{
		RequestID: req.RequestID,
	}

	return result, nil
}
