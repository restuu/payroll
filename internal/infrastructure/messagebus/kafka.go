package messagebus

import (
	"context"
	"log/slog"

	"payroll/internal/app"
	"payroll/internal/infrastructure/config"
	"payroll/internal/infrastructure/log"
	"payroll/internal/presentation/messaging"
	"payroll/internal/presentation/messaging/topics"

	"github.com/twmb/franz-go/pkg/kgo"
)

// NewKafkaClient creates a new franz-go Kafka client.
// It returns the client which should be closed on application shutdown.
func NewKafkaClient(cfg config.KafkaConfig) (*kgo.Client, error) {
	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ConsumerGroup(cfg.ConsumerGroup),
		kgo.ConsumeTopics(getTopicNames(cfg)...),
		kgo.DisableAutoCommit(),
		kgo.AllowAutoTopicCreation(),
		kgo.ConsumeStartOffset(kgo.NewOffset().AtEnd()),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtEnd()),
		// Add other options like authentication, group id, etc. as needed
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.Background()); err != nil {
		return nil, err
	}

	return client, nil
}

func StartConsumer(ctx context.Context, client *kgo.Client, cfg config.KafkaConfig, services *app.Services) {
	slog.InfoContext(ctx, "Start consuming for topics", slog.Any("topics", client.GetConsumeTopics()))

	for {
		fetches := client.PollFetches(ctx)
		if fetches.IsClientClosed() || ctx.Err() != nil {
			slog.InfoContext(ctx, "consumer stopped", "client_closed", fetches.IsClientClosed(), "ctx_err", ctx.Err())
			return
		}

		if errs := fetches.Errors(); len(errs) > 0 {
			// All errors are retried internally, so handle FATAL errors only.
			// We log as an error instead of fatal to prevent a single partition
			// issue from crashing the entire application.
			for _, err := range errs {
				slog.ErrorContext(ctx, "kafka consumer error",
					slog.String("topic", err.Topic),
					slog.Int("partition", int(err.Partition)),
					slog.Any("error", err.Err),
				)
			}
		}

		successfulRecords := make([]*kgo.Record, 0, fetches.NumRecords())

		// Process all records from the fetched batch.
		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			if err := messaging.Dispatch(record, cfg, services); err == nil {
				successfulRecords = append(successfulRecords, record)
			}
		}

		if err := client.CommitRecords(ctx, successfulRecords...); err != nil {
			slog.ErrorContext(ctx, "failed to commit records", log.WithErrorAttr(err))
		}
	}
}

func getTopicNames(cfg config.KafkaConfig) []string {
	topicNames := []string{
		topics.PayrollGenerate,
	}

	if len(cfg.Topics) == 0 {
		return topicNames
	}

	newTopicNames := make([]string, len(topicNames))

	for i, name := range topicNames {
		newTopicNames[i] = cfg.Topics.GetName(name)
	}

	return newTopicNames
}
