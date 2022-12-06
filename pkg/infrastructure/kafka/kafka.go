// Package kafka provide simple abstraction for integrate with kafka
package kafka

import (
	"context"
	"crypto/tls"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
)

var _producer sarama.SyncProducer

// InitSaramaLog init the internal sarama log
func InitSaramaLog() {
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
}

// InitSaramaProducer init kafka producer
func InitSaramaProducer(brokers string, clientID string, username string, password string) {
	_producer = newProducer(brokers, clientID, username, password)
}

// NewConsumerAndSubscribe create new instance of consumer and begin subscribe to topics
func NewConsumerAndSubscribe(brokers string,
	clientID string,
	username string,
	password string,
	consumerGroupID string,
	topics []string,
	kafkaConsumerHandler func(*sarama.ConsumerMessage) error,
) {
	config := createSaramaConsumerConfig(clientID, username, password)
	consumerGroup := createSaramaConsumerGroup(config, brokers, consumerGroupID)
	consumer := createConsumerInstance(kafkaConsumerHandler)
	subscribe(consumerGroup, consumer, topics)
}

func createSaramaConsumerGroup(config *sarama.Config, brokers string, consumerGroupID string) sarama.ConsumerGroup {
	consumerGroup, err := sarama.NewConsumerGroup(getBrokers(brokers), consumerGroupID, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	return consumerGroup
}

func createConsumerInstance(kafkaConsumerHandler func(*sarama.ConsumerMessage) error) *Consumer {
	consumer := &Consumer{
		ready:                make(chan bool),
		kafkaConsumerHandler: kafkaConsumerHandler,
	}
	return consumer
}

func subscribe(consumerGroup sarama.ConsumerGroup, consumer *Consumer, topics []string) {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := consumerGroup.Consume(ctx, topics, consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	keepRunning := true
	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		}
	}

	cancel()
	wg.Wait()

	if err := consumerGroup.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}

func newProducer(brokers string, clientID string, username string, password string) sarama.SyncProducer {
	config := createSaramaProducerConfig(clientID, username, password)

	producer, err := sarama.NewSyncProducer(getBrokers(brokers), config)
	if err != nil {
		log.Panicf("Failed to start Sarama producer: %v", err)
	}

	return producer
}

// Produce do produce a message to a topic
func Produce(topic string, value []byte, key []byte) error {
	if _producer == nil {
		log.Panicf("[kafka.Produce] try to produce before InitSaramaProducer()")
	}

	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(value),
		Key:   sarama.StringEncoder(key),
	}

	partition, offset, err := _producer.SendMessage(message)
	if err != nil {
		log.Printf("[kafka.Produce] error produce message %v\n", err)
		return err
	}

	log.Printf("[kafka.Produce] successfully produce message: partition = %d, offset = %d\n", partition, offset)
	return nil
}

func createSaramaConsumerConfig(clientID string, username string, password string) *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.ClientID = clientID
	config.Net.SASL.Enable = true
	config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	config.Net.SASL.User = username
	config.Net.SASL.Password = password
	config.Net.SASL.Handshake = true
	config.Net.TLS.Enable = true
	config.Net.TLS.Config = &tls.Config{
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
	}
	return config
}

func createSaramaProducerConfig(clientID string, username string, password string) *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.ClientID = clientID
	config.Net.SASL.Enable = true
	config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	config.Net.SASL.User = username
	config.Net.SASL.Password = password
	config.Net.SASL.Handshake = true
	config.Net.TLS.Enable = true
	config.Net.TLS.Config = &tls.Config{
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
	}
	return config
}

func getBrokers(brokers string) []string {
	return strings.Split(brokers, ",")
}
