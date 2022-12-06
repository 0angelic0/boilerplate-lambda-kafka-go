// Package main is a main program
package main

import (
	"log"
	"os"

	"boilerplate-lambda-kafka-go/pkg/demo/event"
	"boilerplate-lambda-kafka-go/pkg/infrastructure/kafka"

	"github.com/Shopify/sarama"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	log.Print("Hello Lambda Kafka Demo")

	// Producer
	initProducer()

	// Consumer
	if os.Getenv("IS_SERVERLESS") == "true" {
		lambda.Start(lambdaHandler)
	} else {
		localHandler()
	}
}

func lambdaHandler(event events.KafkaEvent) error {
	for _, records := range event.Records {
		for i, record := range records {
			log.Printf("[lambdaHandler] will process message: value = %s, timestamp = %v, topic = %s, offset = %d\n", record.Value, record.Timestamp, record.Topic, record.Offset)
			message := kafka.ToMessageFromAwsKafkaRecord(&records[i])
			err := executeHandler(message)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func localHandler() {
	brokers := os.Getenv("KAFKA_BROKERS")
	clientID := os.Getenv("KAFKA_CLIENT_ID")
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")
	consumerGroupID := os.Getenv("KAFKA_CONSUMER_GROUP_ID")
	topics := []string{
		os.Getenv("KAFKA_TOPIC"),
	}

	kafkaConsumerHandler := func(saramaMessage *sarama.ConsumerMessage) error {
		log.Printf("[localHandler] will process message: %v\n", saramaMessage)
		message := kafka.ToMessageFromSaramaConsumerMessage(saramaMessage)
		return executeHandler(message)
	}
	kafka.NewConsumerAndSubscribe(brokers, clientID, username, password, consumerGroupID, topics, kafkaConsumerHandler)
}

func executeHandler(message *kafka.Message) error {
	return event.Handler(message)
}

func initProducer() {
	brokers := os.Getenv("KAFKA_BROKERS")
	clientID := os.Getenv("KAFKA_CLIENT_ID")
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")
	topic := os.Getenv("KAFKA_TOPIC")

	kafka.InitSaramaLog()
	kafka.InitSaramaProducer(brokers, clientID, username, password)

	err := kafka.Produce(topic, []byte("Hello World"), []byte("Hello Key"))
	if err != nil {
		log.Panicf("error produce message: %v", err)
	}
}
