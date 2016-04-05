// Package sns provides a broker implementation that will forward messages to
// an SNS topic.
package sns

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"

	"github.com/lunohq/relay/broker"
	"github.com/lunohq/relay/broker/turnstile"
	log "github.com/Sirupsen/logrus"
)

type Turnstiles []turnstile.Turnstile

// Options are the options the Broker accepts.
type Options struct {
	// TopicArn is the ARN for the topic we want to publish to.
	TopicArn string
	// Turnstiles is an array of broker.Turnstile for the current broker
	Turnstiles Turnstiles
}

// Broker is an implementation of the broker.Broker interface.
type Broker struct {
	broker.Broker
	// TopicArn is the ARN for the topic we want to publish to.
	TopicArn string
	// Turnstiles is an array of broker.Turnstile to pass an event through.
	Turnstiles Turnstiles
}

// New returns a new Broker backed by SNS
func New(options Options) *Broker {
	messageTurnstiles := []turnstile.Turnstile{
		&turnstile.MessageEvents{},
		&turnstile.IgnoreOwnMessages{},
		&turnstile.OnlyMentionOrDM{},
	}

	turnstiles := []turnstile.Turnstile{
		turnstile.NewTurnstileGroup(messageTurnstiles),
		&turnstile.ConnectedEvent{},
	}
	return &Broker{
		TopicArn: options.TopicArn,
		Turnstiles: turnstiles,
	}
}

// ShouldHandle passes the event through any registered turnstiles to determine if
// the event should be forwarded to the broker.
func (b *Broker) ShouldHandle(e broker.Event) bool {
	for _, t := range b.Turnstiles {
		if handle, msg := t.Test(e); !handle {
			log.WithFields(log.Fields{
				"event": e,
			}).Info(msg)
		} else {
			return handle
		}
	}
	return false
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
