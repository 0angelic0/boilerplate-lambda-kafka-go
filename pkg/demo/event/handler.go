// Package event hold handler for processing event bus
package event

import (
	"log"

	"boilerplate-lambda-kafka-go/pkg/infrastructure/kafka"
)

// Handler do handle the messages in event bus
func Handler(message *kafka.Message) error {
	log.Printf("[demo.event.Handler] will process message: value = %s, timestamp = %v, topic = %s, offset = %d\n", message.Value, message.Timestamp, message.Topic, message.Offset)
	return nil
}
