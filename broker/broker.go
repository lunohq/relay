// Package broker provides methods for sending events to various broker backends
package broker

type Context struct {
	// BotID is the user id of the bot
	BotID string `json:"bot_id"`
	// TeamID is the id of the team
	TeamID string `json:"team_id"`
}

type Event struct {
	// Type is the type of event.
	Type string `json:"type"`
	// Source is the source of the event. This allows us to support any relay specific
	// messages we may want to send to the broker backend in the future.
	Source string `json:"source"`
	// Payload is the payload of the event.
	Payload interface{} `json:"payload"`
	// Context is any relevant context for the event.
	Context Context `json:"context"`
}

// Broker represents something that can handle a relay event
type Broker interface {
	// Handle should process an event from the slack client. In general, it's
	// expected that the event will be forwarded to some other system for
	// processing (ie. SNS).
	Handle(e Event) error
}
