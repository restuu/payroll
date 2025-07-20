package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

// examplePayloadA defines the expected JSON structure for messages on "topic-a"
type examplePayloadA struct {
	UserID  string `json:"userId"`
	Event   string `json:"event"`
	Details string `json:"details"`
}

// HandleTopicA processes messages from "topic-a", expecting a JSON payload.
func HandleTopicA(record *kgo.Record) {
	fmt.Printf("Handling Topic A: Partition %d, Offset %d\n", record.Partition, record.Offset)

	var payload examplePayloadA
	if err := json.Unmarshal(record.Value, &payload); err != nil {
		log.Printf("ERROR: Could not unmarshal JSON from topic-a: %v. Raw value: %s", err, string(record.Value))
		return
	}

	// Successfully unmarshaled, now handle the payload
	fmt.Printf("  => Processed Payload A: UserID=%s, Event=%s, Details='%s'\n", payload.UserID, payload.Event, payload.Details)
}
