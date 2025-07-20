package messaging

import (
	"log/slog"

	"payroll/internal/app"
	"payroll/internal/infrastructure/config"
	"payroll/internal/infrastructure/log"
	"payroll/internal/presentation/messaging/handlers"
	"payroll/internal/presentation/messaging/topics"

	"github.com/twmb/franz-go/pkg/kgo"
)

// Dispatch routes the incoming Kafka record to the appropriate handler based on its topic.
func Dispatch(record *kgo.Record, cfg config.KafkaConfig, services *app.Services) (err error) {
	slog.InfoContext(record.Context, "Dispatch message", "topic", record.Topic, "x", cfg.Topics.GetName(topics.PayrollGenerate))

	switch record.Topic {
	case cfg.Topics.GetName(topics.PayrollGenerate):
		err = handlers.HandlePayrollGenerate(record, services.PayrollService)
	default:
		handlers.HandleDefault(record)
	}

	if err != nil {
		slog.ErrorContext(record.Context, "Dispatch message failed",
			slog.String("topic", cfg.Topics.GetName(topics.PayrollGenerate)),
			slog.String("partition", string(record.Partition)),
			slog.Int64("offset", record.Offset),
			log.WithErrorAttr(err),
		)

		return err
	}

	return nil
}
