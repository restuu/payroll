package handlers

import (
	"log/slog"

	"github.com/twmb/franz-go/pkg/kgo"
)

// HandleDefault processes messages from any other topic.
func HandleDefault(record *kgo.Record) {
	slog.Warn("Received message from unhandled topic",
		slog.String("topic", record.Topic),
		slog.String("key", string(record.Key)),
		slog.Int64("offset", record.Offset))
}
