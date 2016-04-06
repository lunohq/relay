// Package lambda provides a handler implementation that will invoke a lambda
// function for each message it receives.
package lambda

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"

	"github.com/lunohq/relay/handler"
	log "github.com/Sirupsen/logrus"
)

// Options are the options the Handler accepts.
type Options struct {
	// FunctionName is the name of the lambda function to invoke.
	FunctionName string
	// Qualifier is the function alias to reference when invoking the function.
	Qualifier string
}

type Handler struct {
	handler.Base

	// FunctionName is the name of the lambda function to invoke.
	FunctionName string
	// Qualifier is the function alias to reference when invoking the function.
	Qualifier string
}

// New returns a new Lambda Handler
func New(options Options) *Handler {
	return &Handler{
		Base: handler.Base{
			Name: "lambda",
		},
		FunctionName: options.FunctionName,
		Qualifier: options.Qualifier,
	}
}

// Handle invokes the given lambda function with the event as the payload.
func (h *Handler) Handle(e handler.Event) error {
	bytes, err := e.Serialize()
	if err == nil {
		log.WithFields(log.Fields{
			"event": e,
			"function_name": h.FunctionName,
		}).Info("invoking lambda function")

		svc := lambda.New(session.New())
		params := &lambda.InvokeInput{
			FunctionName: aws.String(h.FunctionName),
			InvocationType: aws.String(lambda.InvocationTypeEvent),
			Payload: bytes,
			Qualifier: aws.String(h.Qualifier),
		}
		_, err = svc.Invoke(params)
	}
	return err
}
