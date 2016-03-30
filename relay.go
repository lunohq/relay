package relay

import (
	log "github.com/Sirupsen/logrus"

	"github.com/lunohq/relay/slack"
)

type Config struct {
	TeamId string
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
		"team_id": r.config.TeamId,
	}).Info("starting relay")

	client := slack.New(r.config.TeamId, r.config.Token)
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
