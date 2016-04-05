// Package firehose provides a handler implementation that will forward messages
// to a Kinesis Firehose Delivery Streamm.
package firehose

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"

	"github.com/lunohq/relay/handler"
	log "github.com/Sirupsen/logrus"
)

// Options are the options the Handler accepts.
type Options struct {
	// DeliveryStreamName is the name of an existing delivery stream
	DeliveryStreamName string
}

type Handler struct {
	handler.Base

	// DeliveryStreamName is the name of an existing delivery stream
	DeliveryStreamName string
}

// New returns a new Firehose Handler
func New(options Options) *Handler {
	return &Handler{
		DeliveryStreamName: options.DeliveryStreamName,
	}
}

// Handle sends the event to the given Firehose Delivery Stream
func (h *Handler) Handle(e handler.Event) error {
	bytes, err := e.Serialize()
	if err == nil {
		log.WithFields(log.Fields{
			"event": e,
		}).Info("forwarding to firehose")

		svc := firehose.New(session.New())
		params := &firehose.PutRecordInput{
			DeliveryStreamName: aws.String(h.DeliveryStreamName),
			Record: &firehose.Record{
				Data: bytes,
			},
		}
		_, err = svc.PutRecord(params)
	}
	return err
}
