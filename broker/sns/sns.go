// Package sns provides a broker implementation that will forward messages to
// an SNS topic.
package sns

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"

	"github.com/lunohq/relay/broker"
	log "github.com/Sirupsen/logrus"
)

// Options are the options the Broker accepts.
type Options struct {
	// TopicArn is the ARN for the topic we want to publish to.
	TopicArn string
}

// Broker is an implementation of the broker.Broker interface.
type Broker struct {
	// TopicArn is the ARN for the topic we want to publish to.
	TopicArn string
}

// New returns a new Broker backed by SNS
func New(options Options) *Broker {
	return &Broker{
		TopicArn: options.TopicArn,
	}
}

// Handle sends the event to the given SNS topic
func (b *Broker) Handle(e broker.Event) error {
	bytes, err := json.Marshal(e)
	if err == nil {
		log.WithFields(log.Fields{
			"event": e,
		}).Info("forwarding to sns")
	}

	svc := sns.New(session.New())
	params := &sns.PublishInput{
		Message: aws.String(string(bytes)),
		TopicArn: aws.String(b.TopicArn),
	}
	_, err = svc.Publish(params)
	return err
}
