package relay

import (
	log "github.com/Sirupsen/logrus"

	"github.com/lunohq/relay/handler"
	"github.com/lunohq/relay/slack"
)

type Config struct {
	// Handlers to use to handle events
	Handlers []handler.Handler

	// TeamID is the id of the team to init the client for.
	TeamID string
	// Token is the token to init the client with.
	Token string
}

type Clients []*slack.Client

type Relay struct {
	config *Config
	clients Clients
}

// New returns a new Relay instance
func New(c *Config) *Relay {
	return &Relay{config: c}
}

// Start opens a RTM connection for the configured team
func (r *Relay) Start() {
	log.WithFields(log.Fields{
		"team_id": r.config.TeamID,
	}).Info("starting relay")

	client := slack.New(slack.Options{
		Handlers: r.config.Handlers,
		TeamID: r.config.TeamID,
		Token: r.config.Token,
	})
	r.clients = append(r.clients, client)

	for _, c := range r.clients {
		c.Connect()
		go c.Start()
	}
}

// Shutdown closes a RTM connection for the configured team
func (r *Relay) Shutdown() error {
	log.Info("shutting down relay")
	for _, c := range r.clients {
		c.Disconnect()
	}

	return nil
}
