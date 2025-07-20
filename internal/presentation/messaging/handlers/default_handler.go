package handlers

import (
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

// HandleDefault processes messages from any other topic.
func HandleDefault(record *kgo.Record) {
	log.Printf("Received message from unhandled topic '%s'. Raw value: %s", record.Topic, string(record.Value))
}
