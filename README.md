# relay

Slack bot message broker that publishes events to SNS. This allows us to easily
consume these messages with AWS Lambda functions.

Roadmap:
v1
- [x] cmd to run relay
- [x] relay connects to Slack's RTM api for a specific slackbot token and maintains the websocket connection
- [x] for each event that is received by relay, we'll foward the message to a specific SNS topic
v2
- [ ] relay should support initializing multiple clients
- [ ] relay should support fetching clients that need to be connected from a data source (DynamoDB or Redis)
- [ ] relay should support listening for new clients to connect
- [ ] relay should support some sort of high availability mode
