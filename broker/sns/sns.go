// Package sns provides a broker implementation that will forward messages to
// an SNS topic.
package sns

import (
	//"encoding/json"

	"github.com/lunohq/relay/broker"
	log "github.com/Sirupsen/logrus"
)

// Broker is an implementation of the broker.Broker interface.
type Broker struct {}

// New returns a new Broker backed by SNS
func New() *Broker {
	return &Broker{}
}

// Handle sends the event to the given SNS topic
func (b *Broker) Handle(e broker.Event) error {
	//bytes, err := json.Marshal(e)
	//if err == nil {
	log.WithFields(log.Fields{
		"event": e,
	}).Info("forwarding to sns")
	//}

	return nil
}
