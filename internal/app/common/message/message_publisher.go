package message

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"payroll/internal/infrastructure/config"
	"payroll/internal/infrastructure/log"

	"github.com/google/uuid"
	"github.com/twmb/franz-go/pkg/kgo"
)

type MessagePublisher interface {
	PublishJSON(ctx context.Context, topic string, payload any) (err error)
}

func NewMessagePublisher(cfg config.KafkaConfig, client *kgo.Client) MessagePublisher {
	return &messagePublisher{
		cfg:    cfg,
		client: client,
	}
}

type messagePublisher struct {
	cfg    config.KafkaConfig
	client *kgo.Client
}

func (m *messagePublisher) PublishJSON(ctx context.Context, topic string, payload any) (err error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("MessagePublisher.PublishJSON failed marshal json: %w", err)
	}

	publishCtx := context.WithoutCancel(ctx)

	topic = m.cfg.Topics.GetName(topic)

	record := &kgo.Record{
		Key:       []byte(uuid.New().String()),
		Topic:     topic,
		Timestamp: time.Now(),
		Value:     b,
		Context:   publishCtx,
	}

	m.client.Produce(publishCtx, record, func(_ *kgo.Record, err error) {
		if err != nil {
			slog.ErrorContext(ctx, "MessagePublisher.PublishJSON failed on client.Produce",
				log.WithErrorAttr(err),
				log.WithMetadata(map[string]any{
					"topic": topic,
				}))
		}
	})

	return nil
}
