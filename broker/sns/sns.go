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

type Filters []broker.Filter

// Options are the options the Broker accepts.
type Options struct {
	// TopicArn is the ARN for the topic we want to publish to.
	TopicArn string
	// Filters is an array of filters for the current broker
	Filters Filters
}

// Broker is an implementation of the broker.Broker interface.
type Broker struct {
	broker.Broker
	// TopicArn is the ARN for the topic we want to publish to.
	TopicArn string
	// Filters is an array of broker.Filter to pass an event through.
	Filters Filters
}

// New returns a new Broker backed by SNS
func New(options Options) *Broker {
	filters := []broker.Filter{
		&broker.OnlyMessageEventsFilter{},
		&broker.IgnoreOwnMessagesFilter{},
		&broker.OnlyMentionsOrDMFilter{},
	}
	return &Broker{
		TopicArn: options.TopicArn,
		Filters: options.Filters,
	}
}

// ShouldHandle passes the event through any registered filters to determine if
// the event should be forwarded to the broker.
func (b *Broker) ShouldHandle(e broker.Event) bool {
	for _, f := range b.Filters {
		if handle, msg := f.ShouldHandle(e); !handle {
			log.WithFields(log.Fields{
				"event": e,
			}).Info(msg)
			return handle
		}
	}
	return true
}

// Handle sends the event to the given SNS topic
func (b *Broker) Handle(e broker.Event) error {
	if !b.ShouldHandle(e) {
		return nil
	}

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
