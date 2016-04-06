package slack

import (
	api "github.com/nlopes/slack"
	log "github.com/Sirupsen/logrus"

	"github.com/lunohq/relay/handler"
)

type Options struct {
	// Handlers to use to handle events
	Handlers []handler.Handler

	// TeamID is the id of the team to init the client for.
	TeamID string
	// Token is the token to init the client with.
	Token string
}

type Client struct {
	// Handlers that will handle events
	Handlers []handler.Handler

	// TeamID is the id of the team for this client.
	TeamID string
	// Token is the token for this client.
	Token string

	// rtm is a pointer to the RTM connection
	rtm *api.RTM

}

// New returns a new Client instance
func New(options Options) *Client {
	log.WithFields(log.Fields{
		"handlers": len(options.Handlers),
		"team_id": options.TeamID,
	}).Info("creating slack client")
	return &Client{
		Handlers: options.Handlers,
		TeamID: options.TeamID,
		Token: options.Token,
	}
}

// Connect connects the client to the RTM api
func (c *Client) Connect() {
	log.WithFields(log.Fields{
		"team_id": c.TeamID,
	}).Info("connecting to slack rtm")

	a := api.New(c.Token)
	a.SetDebug(true)

	c.rtm = a.NewRTM()
	go c.rtm.ManageConnection()
}

// Start listens to events from the RTM api
func (c *Client) Start() error {
	log.WithFields(log.Fields{
		"team_id": c.TeamID,
	}).Info("starting to listen for events")

	for {
		select {
		case msg := <-c.rtm.IncomingEvents:
			c.Forward(msg)
		}
	}

	return nil
}

// Disconnect disconnects the client from the RTM api
func (c *Client) Disconnect() error {
	log.WithFields(log.Fields{
		"team_id": c.TeamID,
	}).Info("disconnecting from rtm")

	return c.rtm.Disconnect()
}

// Foward fowards an RTMEvent to any handlers
func (c *Client) Forward(e api.RTMEvent) {
	info := c.rtm.GetInfo()
	if info != nil {
		event := handler.Event{
			Type: e.Type,
			// TODO this should be a constant
			Source: "slack",
			Payload: e.Data,
			Context: handler.Context{
				BotID: info.User.ID,
				TeamID: c.TeamID,
			},
			RTMEvent: &e,
		}
		for _, h := range c.Handlers {
			log.WithFields(log.Fields{
				"handler": h.String(),
			}).Info("processing event")
			err := handler.Process(h, event)
			if err != nil {
				log.WithFields(log.Fields{
					"team_id": c.TeamID,
					"err": err,
					"event": event,
				}).Error("Handler failed to process event")
			}
		}
	}
}
