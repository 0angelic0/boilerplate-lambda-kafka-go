// Package main is a main program
package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"boilerplate-lambda-kafka-go/pkg/demo/event"
	"boilerplate-lambda-kafka-go/pkg/infrastructure/kafka"

	"github.com/Shopify/sarama"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/riferrei/srclient"
)

func main() {
	log.Print("Hello Lambda Kafka Demo-Schema")

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
	// 2) Create a instance of the client to retrieve the schemas for each message
	schemaRegistryClient := srclient.CreateSchemaRegistryClient(os.Getenv("SCHEMA_REGISTRY_ENDPOINT"))
	schemaRegistryClient.SetCredentials(os.Getenv("SCHEMA_REGISTRY_USERNAME"), os.Getenv("SCHEMA_REGISTRY_PASSWORD"))

	// 3) Recover the schema id from the message and use the
	// client to retrieve the schema from Schema Registry.
	// Then use it to deserialize the record accordingly.
	schemaID := binary.BigEndian.Uint32(message.Value[1:5])
	log.Printf("schemaID = %d", schemaID)
	schema, err := schemaRegistryClient.GetSchema(int(schemaID))
	if err != nil {
		panic(fmt.Sprintf("Error getting the schema with id '%d' %s", schemaID, err))
	}
	native, _, _ := schema.Codec().NativeFromBinary(message.Value[5:])
	value, _ := schema.Codec().TextualFromNative(nil, native)
	fmt.Printf("Here is the message %s\n", string(value))

	message.Value = value

	return event.Handler(message)
}

func initProducer() {
	brokers := os.Getenv("KAFKA_BROKERS")
	clientID := os.Getenv("KAFKA_CLIENT_ID")
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")
	topic := os.Getenv("KAFKA_TOPIC")

	type TopicMessageBody struct {
		MyField1 int32   `json:"my_field1"`
		MyField2 float64 `json:"my_field2"`
		MyField3 string  `json:"my_field3"`
	}

	kafka.InitSaramaLog()
	kafka.InitSaramaProducer(brokers, clientID, username, password)

	// 2) Fetch the latest version of the schema, or create a new one if it is the first
	schemaRegistryClient := srclient.CreateSchemaRegistryClient(os.Getenv("SCHEMA_REGISTRY_ENDPOINT"))
	schemaRegistryClient.SetCredentials(os.Getenv("SCHEMA_REGISTRY_USERNAME"), os.Getenv("SCHEMA_REGISTRY_PASSWORD"))

	schema, err := schemaRegistryClient.GetLatestSchema(os.Getenv("SCHEMA_REGISTRY_SUBJECT"))
	if schema == nil {
		// schemaBytes, _ := ioutil.ReadFile("complexType.avsc")
		// schema, err = schemaRegistryClient.CreateSchema(topic, string(schemaBytes), srclient.Avro)
		// if err != nil {
		// 	panic(fmt.Sprintf("Error creating the schema %s", err))
		// }
		panic(fmt.Sprintf("Error reading the schema %s", err))
	}
	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(schema.ID()))

	// 3) Serialize the record using the schema provided by the client,
	// making sure to include the schema id as part of the record.
	newComplexType := TopicMessageBody{MyField1: 1, MyField2: 0.028, MyField3: "HeyHey"}
	value, _ := json.Marshal(newComplexType)
	native, _, _ := schema.Codec().NativeFromTextual(value)
	valueBytes, _ := schema.Codec().BinaryFromNative(nil, native)

	var recordValue []byte
	recordValue = append(recordValue, byte(0))
	recordValue = append(recordValue, schemaIDBytes...)
	recordValue = append(recordValue, valueBytes...)

	err = kafka.Produce(topic, recordValue, []byte("Hello Key"))
	if err != nil {
		log.Panicf("error produce message: %v", err)
	}
}
