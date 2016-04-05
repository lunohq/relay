package handler

import (
	"strings"

	"github.com/nlopes/slack"
)

type Turnstile interface {
	// Test should take a Event and return whether or not the event
	// should be handled by the handler.
	Test(e Event) (bool, string)
}

// TurnstileGroup can be used to operate on a group of Turnstiles with either
// ALL or ANY logic
type TurnstileGroup struct {
	all bool

	Turnstiles []Turnstile
}

func NewTurnstileGroup(turnstiles []Turnstile, all bool) *TurnstileGroup {
	return &TurnstileGroup{Turnstiles: turnstiles, all: all}
}

func (g *TurnstileGroup) Test(e Event) (bool, string) {
	for _, t := range g.Turnstiles {
		handle, msg := t.Test(e)
		if g.all && !handle {
			return handle, msg
		} else if !g.all && handle {
			return handle, msg
		}
	}

	if g.all {
		return true, ""
	} else {
		return false, ""
	}
}

// AllowMessages allows slack.MessageEvent
type AllowMessages struct {}

func (t *AllowMessages) Test(e Event) (bool, string) {
	switch e.RTMEvent.Data.(type) {
	case *slack.MessageEvent:
		return true, ""
	default:
		return false, "filtering out non slack.MessageEvent"
	}
}

// IgnoreOwnMessages filters out any messages sent by the connected bot
// (avoids handling messages we've sent)
type IgnoreOwnMessages struct {}

func (t *IgnoreOwnMessages) Test(e Event) (bool, string) {
	ev, ok := e.RTMEvent.Data.(*slack.MessageEvent)
	if ok && ev.User == e.Context.BotID {
		return false, "filtering out message sent by current bot"
	}
	return true, ""
}

// AllowMention allows an mention of the connected bot
type AllowMention struct {}

func (t *AllowMention) Test(e Event) (bool, string) {
	ev, ok := e.RTMEvent.Data.(*slack.MessageEvent)
	if ok {
		if strings.Contains(ev.Text, e.Context.BotID) {
			return true, ""
		}
	}
	return false, "ignoring non mention of connected bot"
}

// AllowDirectMessage allows a direct message to the connected bot
type AllowDirectMessage struct {}

func (t *AllowDirectMessage) Test(e Event) (bool, string) {
	ev, ok := e.RTMEvent.Data.(*slack.MessageEvent)
	if ok {
		if strings.HasPrefix(ev.Channel, "D") {
			return true, ""
		}
	}
	return false, "ignoring non direct message of connected bot"
}

// AllowConnectedEvent allows slack.ConnectedEvent
type AllowConnectedEvent struct {}

func (t *AllowConnectedEvent) Test(e Event) (bool, string) {
	_, ok := e.RTMEvent.Data.(*slack.ConnectedEvent)
	if ok {
		return true, ""
	}
	return false, "ignoring non slack.ConnectedEvent"
}
