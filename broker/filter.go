package broker

import (
	"strings"

	"github.com/nlopes/slack"
)

type Filter interface {
	// ShouldHandle should take a broker Event and return whether or not the
	// event should be handled by the broker.
	ShouldHandle(e Event) (bool, string)
}

// OnlyMessageEventsFilter filters out anything but the slack.MessageEvent
type OnlyMessageEventsFilter struct {
	Filter
}

func (f *OnlyMessageEventsFilter) ShouldHandle(e Event) (bool, string) {
	switch e.RTMEvent.Data.(type) {
	case *slack.MessageEvent:
		return true, ""
	default:
		return false, "filtering out non-message event"
	}
}

// IgnoreOwnMessagesFilter filters out any messages sent by the connected bot
// (avoids handling messages we've sent)
type IgnoreOwnMessagesFilter struct {
	Filter
}

func (f *IgnoreOwnMessagesFilter) ShouldHandle(e Event) (bool, string) {
	ev, ok := e.RTMEvent.Data.(*slack.MessageEvent)
	if ok && ev.User == e.Context.BotID {
		return false, "filtering out message sent by current bot"
	}
	return true, ""
}

// OnlyMentionsOrDMFilter filters out anything but a direct message to or mention of the connected bot
type OnlyMentionsOrDMFilter struct {
	Filter
}

func (f *OnlyMentionsOrDMFilter) ShouldHandle(e Event) (bool, string) {
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
