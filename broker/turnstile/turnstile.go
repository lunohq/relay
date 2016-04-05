package turnstile

import (
	"strings"

	"github.com/nlopes/slack"
	"github.com/lunohq/relay/broker"
)

type Turnstile interface {
	// Turnstile should take a broker Event and return whether or not the event
	// should be handled by the broker.
	Test(e broker.Event) (bool, string)
}

type TurnstileGroup struct {
	Turnstiles []Turnstile
}

func NewTurnstileGroup(turnstiles []Turnstile) *TurnstileGroup {
	return &TurnstileGroup{Turnstiles: turnstiles}
}

func (t *TurnstileGroup) Test(e broker.Event) (bool, string) {
	for _, t := range t.Turnstiles {
		if handle, msg := t.Test(e); !handle {
			return handle, msg
		}
	}
	return true, ""
}

// MessageEvents allows slack.MessageEvent
type MessageEvents struct {}

func (t *MessageEvents) Test(e broker.Event) (bool, string) {
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

func (t *IgnoreOwnMessages) Test(e broker.Event) (bool, string) {
	ev, ok := e.RTMEvent.Data.(*slack.MessageEvent)
	if ok && ev.User == e.Context.BotID {
		return false, "filtering out message sent by current bot"
	}
	return true, ""
}

// OnlyMentionOrDM filters out anything but a direct message to or mention of the connected bot
type OnlyMentionOrDM struct {}

func (t *OnlyMentionOrDM) Test(e broker.Event) (bool, string) {
	ev, ok := e.RTMEvent.Data.(*slack.MessageEvent)
	if ok {
		if ev.Channel == e.Context.BotID {
			return true, ""
		} else if strings.Contains(ev.Text, e.Context.BotID) {
			return true, ""
		}
	}
	return false, "filtering out non DM or mention of connected bot"
}

// ConnectedEvent allows slack.ConnectedEvent
type ConnectedEvent struct {}

func (t *ConnectedEvent) Test(e broker.Event) (bool, string) {
	_, ok := e.RTMEvent.Data.(*slack.ConnectedEvent)
	if ok {
		return true, ""
	}
	return false, "filtering out non slack.ConnectedEvent"
}
