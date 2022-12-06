package kafka

import (
	"encoding/base64"
	"log"

	"github.com/Shopify/sarama"
	"github.com/aws/aws-lambda-go/events"
)

// ToMessageFromAwsKafkaRecord convert KafkaRecord to Message
func ToMessageFromAwsKafkaRecord(awsKafkaRecord *events.KafkaRecord) *Message {
	key, err := base64.StdEncoding.DecodeString(awsKafkaRecord.Key)
	if err != nil {
		log.Panicf("[ToMessageFromAwsKafkaRecord] error decoding message key: %v", err)
	}

	value, err := base64.StdEncoding.DecodeString(awsKafkaRecord.Value)
	if err != nil {
		log.Panicf("[ToMessageFromAwsKafkaRecord] error decoding message value: %v", err)
	}

	message := &Message{
		Topic:     awsKafkaRecord.Topic,
		Partition: awsKafkaRecord.Partition,
		Offset:    awsKafkaRecord.Offset,
		// Headers:   awsKafkaRecord.Headers,
		Headers:   nil,
		Key:       key,
		Value:     value,
		Timestamp: awsKafkaRecord.Timestamp.Time,
	}
	return message
}

// ToMessageFromSaramaConsumerMessage convert ConsumerMessage to Message
func ToMessageFromSaramaConsumerMessage(saramaMessage *sarama.ConsumerMessage) *Message {
	message := &Message{
		Topic:     saramaMessage.Topic,
		Partition: int64(saramaMessage.Partition),
		Offset:    saramaMessage.Offset,
		// Headers:   saramaMessage.Headers,
		Headers:   nil,
		Key:       saramaMessage.Key,
		Value:     saramaMessage.Value,
		Timestamp: saramaMessage.Timestamp,
	}
	return message
}
