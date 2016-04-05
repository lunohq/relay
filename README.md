# relay

Slack bot message handler that publishes events to SNS. This allows us to easily
consume these messages with AWS Lambda functions.

## Roadmap

### v1
- [x] cmd to run relay
- [x] relay connects to Slack's RTM api for a specific slackbot token and maintains the websocket connection
- [x] for each event that is received by relay, we'll foward the message to a specific SNS topic

### v2
- [ ] relay should support initializing multiple clients
- [ ] relay should support fetching clients that need to be connected from a data source (DynamoDB or Redis)
- [ ] relay should support listening for new clients to connect
- [ ] relay should support some sort of high availability mode

## Architecture

- Relay handles connecting Slack clients
- The Slack client maintains the RTM connection with the API, any event it
  receives, it passes to one or many registered Handlers
- The Handler is responsible for handling the event
    - a Handler can have one or many turnstiles registered.
    - A Turnstile controls whether or not the event should be handled by the Handler

Handlers:
- SNS: forward messages to an SNS topic
    - topic ARN
- Lambda: invoke a lambda function directly
    - function ARN
- Kinesis: write the event to a kinesis stream directly
    - can optionally configure a kinesis firehose to write records to an s3 bucket

Our use case:
- Lambda handler invokes a lambda function for each event it receives
- Kinesis handler writes each event to a kinesis stream


We still need something like Redis to prevent multiple handlers from being
invoked for the same message across processes
