// Package handler provides methods for sending events to various handler backends
package handler

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"

	"github.com/nlopes/slack"
)

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
	// messages we may want to send to the handler backend in the future.
	Source string `json:"source"`
	// Payload is the payload of the event.
	Payload interface{} `json:"payload"`
	// Context is any relevant context for the event.
	Context Context `json:"context"`
	// RTMEvent is the raw real time message from slack.
	RTMEvent *slack.RTMEvent `json:"-"`
}

func (e *Event) Serialize() ([]byte, error) {
	bytes, err := json.Marshal(e)
	if err != nil {
		log.WithFields(log.Fields{
			"event": e,
		}).Error("failed to serialize eveent")
	}
	return bytes, err
}

// Handler represents something that can handle a relay event
type Handler interface {

	// Handle should handle an event from the slack client. In general, it's
	// expected that the event will be forwarded to some other system for
	// processing (ie. SNS).
	Handle(e Event) error

	// Process should process an event from the slack client. In general, it's
	// expected that the handler will determine if the event should be handled
	// and then handle it.
	Process(e Event) error

	// Return the name of the handler
	String() string
}

// Base is the base struct that Handlers should embed
type Base struct {
	Handler

	// Name is the name of the handler
	Name string

	// Turnstiles is an array of Turnstile to pass an event through
	Turnstiles []Turnstile
}

// Handle should be overridden by subclass
func (h *Base) Handle(e Event) error {
	return nil
}

// String should return the name of the handler
func (h *Base) String() string {
	return h.Name
}

// ShouldHandle passes the event through any registered turnstiles to determine if
// the event should be handled by the handler.
func (h *Base) ShouldHandle(e Event) bool {
	for _, t := range h.Turnstiles {
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

// Process processes an event
func (h *Base) Process(e Event) error {
	if !h.ShouldHandle(e) {
		return nil
	}

	return h.Handle(e)
}
