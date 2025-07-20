package messaging

import (
	"payroll/internal/presentation/messaging/handlers"

	"github.com/twmb/franz-go/pkg/kgo"
)

// Dispatch routes the incoming Kafka record to the appropriate handler based on its topic.
func Dispatch(record *kgo.Record) {
	var err error
	defer func() {
		if err != nil {

		}
	}()
	switch record.Topic {
	case "topic-a":
		handlers.HandleTopicA(record)
	case "topic-b":
		handlers.HandleTopicB(record)
	default:
		handlers.HandleDefault(record)
	}
}
