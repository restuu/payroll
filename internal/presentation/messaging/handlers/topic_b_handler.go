package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

// examplePayloadB defines the expected JSON structure for messages on "topic-b"
type examplePayloadB struct {
	TransactionID string  `json:"transactionId"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
}

// HandleTopicB processes messages from "topic-b", expecting a JSON payload.
func HandleTopicB(record *kgo.Record) {
	fmt.Printf("Handling Topic B: Partition %d, Offset %d\n", record.Partition, record.Offset)

	var payload examplePayloadB
	if err := json.Unmarshal(record.Value, &payload); err != nil {
		log.Printf("ERROR: Could not unmarshal JSON from topic-b: %v. Raw value: %s", err, string(record.Value))
		return
	}

	// Successfully unmarshaled, now handle the payload
	fmt.Printf("  => Processed Payload B: TransactionID=%s, Amount=%.2f, Currency=%s\n", payload.TransactionID, payload.Amount, payload.Currency)
}
