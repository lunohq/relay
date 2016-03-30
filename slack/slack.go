package slack

import (
	api "github.com/nlopes/slack"
	log "github.com/Sirupsen/logrus"
)

type Client struct {
	TeamId string
	Token string

	rtm *api.RTM
}

// New returns a new Client instance
func New(teamId, token string) *Client {
	return &Client{TeamId: teamId, Token: token}
}

// Connect connects the client to the RTM api
func (c *Client) Connect() {
	log.WithFields(log.Fields{
		"team_id": c.TeamId,
	}).Info("connecting to slack rtm")

	a := api.New(c.Token)
	a.SetDebug(true)

	c.rtm = a.NewRTM()
	go c.rtm.ManageConnection()
}

// Start listens to events from the RTM api
func (c *Client) Start() error {
	log.WithFields(log.Fields{
		"team_id": c.TeamId,
	}).Info("starting to listen for events")

	for {
		select {
		case msg := <-c.rtm.IncomingEvents:
			switch msg.Data.(type) {
			default:
				//log.WithFields(log.Fields{
					//"team_id": c.TeamId,
					//"data": msg.Data,
				//}).Info("Received event")
			}
		}
	}

	return nil
}

// Disconnect disconnects the client from the RTM api
func (c *Client) Disconnect() error {
	log.WithFields(log.Fields{
		"team_id": c.TeamId,
	}).Info("disconnecting from rtm")

	return c.rtm.Disconnect()
}
