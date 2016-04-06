// Package sns provides a handler implementation that will forward messages to
// an SNS topic.
package sns

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"

	"github.com/lunohq/relay/handler"
	log "github.com/Sirupsen/logrus"
)

// Options are the options the Handler accepts.
type Options struct {
	// TopicArn is the ARN for the topic we want to publish to.
	TopicArn string
}

// Handler is an implementation of the handler.Handler interface.
type Handler struct {
	handler.Base

	// TopicArn is the ARN for the topic we want to publish to.
	TopicArn string
}

// New returns a new SNS Handler
func New(options Options) *Handler {
	messageTurnstiles := []handler.Turnstile{
		&handler.AllowMessages{},
		&handler.IgnoreOwnMessages{},
		handler.NewTurnstileGroup([]handler.Turnstile{
			&handler.AllowMention{},
			&handler.AllowDirectMessage{},
		}, false),
	}

	turnstiles := []handler.Turnstile{
		handler.NewTurnstileGroup(messageTurnstiles, true),
		&handler.AllowConnectedEvent{},
	}
	return &Handler{
		Base: handler.Base{
			Name: "sns",
			Turnstiles: turnstiles,
		},
		TopicArn: options.TopicArn,
	}
}

// Handle sends the event to the given SNS topic
func (h *Handler) Handle(e handler.Event) error {

	bytes, err := e.Serialize()
	if err == nil {
		log.WithFields(log.Fields{
			"event": e,
		}).Info("forwarding to sns")

		svc := sns.New(session.New())
		params := &sns.PublishInput{
			Message: aws.String(string(bytes)),
			TopicArn: aws.String(h.TopicArn),
		}
		_, err = svc.Publish(params)
	}
	return err
}
