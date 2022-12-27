#!/bin/bash

export KAFKA_BROKERS=pkc-1dkx6.ap-southeast-1.aws.confluent.cloud:9092
export KAFKA_CLIENT_ID=local-boilerplate-lambda-kafka-go
export KAFKA_CONSUMER_GROUP_ID=local-boilerplate-lambda-kafka-go-consumer-group
export KAFKA_USERNAME=
export KAFKA_PASSWORD=
export KAFKA_TOPIC=dev_king-test-topic-with-schema
export SCHEMA_REGISTRY_ENDPOINT=https://psrc-znpo0.ap-southeast-2.aws.confluent.cloud
export SCHEMA_REGISTRY_USERNAME=
export SCHEMA_REGISTRY_PASSWORD=
export SCHEMA_REGISTRY_SUBJECT=dev_king-test-topic-with-schema-value

export IS_OFFLINE=true 

go run ./services/demo-schema/.