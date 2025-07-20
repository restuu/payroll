package messagebus

import (
	"context"
	"log"
	"log/slog"

	"payroll/internal/infrastructure/config"
	"payroll/internal/presentation/messaging"

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

func StartConsumer(ctx context.Context, client *kgo.Client) {
	for {
		fetches := client.PollFetches(ctx)
		if fetches.IsClientClosed() || ctx.Err() != nil {
			return
		}

		if errs := fetches.Errors(); len(errs) > 0 {
			// All errors are retried internally, so handle FATAL errors only.
			for _, err := range errs {
				log.Fatalf("FATAL error consuming from topic %s: %v", err.Topic, err.Err)
			}
		}

		if fetches.NumRecords() == 0 {
			continue
		}

		iter := fetches.RecordIter()
		if iter.Done() {
			continue
		}
		record := iter.Next()
		messaging.Dispatch(record)
		if err := client.CommitRecords(record.Context, record); err != nil {
			slog.ErrorContext(
				ctx, "Failed to commit record",
				slog.String("topic", record.Topic),
				slog.Int64("partition", int64(record.Partition)),
				slog.Int64("offset", record.Offset),
				slog.Any("error", err),
			)
		}
	}
}

func getTopicNames(cfg config.KafkaConfig) (topics []string) {
	for k, topicConfig := range cfg.Topics {
		name := topicConfig.Name
		if name == "" {
			name = k
		}
		topics = append(topics, name)
	}

	return topics
}
