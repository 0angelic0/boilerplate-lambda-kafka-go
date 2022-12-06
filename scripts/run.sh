#!/bin/bash

export KAFKA_BROKERS=pkc-1dkx6.ap-southeast-1.aws.confluent.cloud:9092
export KAFKA_CLIENT_ID=local-boilerplate-lambda-kafka-go
export KAFKA_CONSUMER_GROUP_ID=local-boilerplate-lambda-kafka-go-consumer-group
export KAFKA_USERNAME=
export KAFKA_PASSWORD=
export KAFKA_TOPIC=dev_king-test-topic

export MYSQL_HOST=xx
export MYSQL_PORT=3306
export MYSQL_USERNAME=xx
export MYSQL_PASSWORD=xx
export MYSQL_DATABASE=xx

export IS_OFFLINE=true 

go run ./services/demo/.