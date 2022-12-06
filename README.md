# boilerplate-lambda-kafka-go

To demo how to setup a project for serverless/lambda (golang) that will be triggered by self-managed kafka event.

The setup purpose is to let developers able to run the project on local (will use sarama as a kafka consumer) and when deploy into lambda environment it seamlessly switch to be triggered by aws lambda-kafka-event (not sarama consumer).

1. run the project on local

We use sarama as a kafka consumer so we need to config kafka broker, username, password, consumer group id, client id to initialize sarama kafka consumer. So, we need to read these configs from ENV.

2. run in lambda environment

We use [lambda-kafka-event](https://docs.aws.amazon.com/lambda/latest/dg/with-kafka.html) via serverless framework [events-kafka](https://www.serverless.com/framework/docs/providers/aws/events/kafka). So, we need to setup `AWS Secret Manager ARN` that store kafka credential (username and password).

We just only need customized starting point for consumer side because there's a different on how we consume the events on local (our application consume from kafka directly) and on lambda (let aws lambda consume from kafka and invoke our function).

## For Producer Side

For producer side, it is same for both on local and on lambda. Our application connect to kafka broker and produce the events directly via sarama's SyncProducer.

## Make commands

To run on local

`make run`

To build for serverless lambda

`make build`

To deploy to dev environment (need aws iam credential client id and secret already setup in `~/.aws/credentials` file)

`make deploy-dev`