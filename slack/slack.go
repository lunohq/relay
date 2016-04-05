package slack

import (
	api "github.com/nlopes/slack"
	log "github.com/Sirupsen/logrus"

	"github.com/lunohq/relay/handler"
)

type Options struct {
	// Handler to use to handle events
	Handler handler.Handler

	// TeamID is the id of the team to init the client for.
	TeamID string
	// Token is the token to init the client with.
	Token string
}

type Client struct {
	// Handler that will handle events
	Handler handler.Handler

	// TeamID is the id of the team for this client.
	TeamID string
	// Token is the token for this client.
	Token string

	// rtm is a pointer to the RTM connection
	rtm *api.RTM

}

// New returns a new Client instance
func New(options Options) *Client {
	return &Client{
		Handler: options.Handler,
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

// Foward fowards an RTMEvent to the handler
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
		err := c.Handler.Process(event)
		if err != nil {
			log.WithFields(log.Fields{
				"team_id": c.TeamID,
				"err": err,
				"event": event,
			}).Error("Handler failed to process event")
		}
	}
}
