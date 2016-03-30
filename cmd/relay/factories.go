package main

import (
	"github.com/lunohq/relay"
	"github.com/lunohq/relay/broker"
	"github.com/lunohq/relay/broker/sns"
	"github.com/codegangsta/cli"
)

func newConfig(c *cli.Context) *relay.Config {
	return &relay.Config{
		Broker: newBroker(c),
		TeamID: c.Args()[0],
		Token: c.Args()[1],
	}
}

func newRelay(c *cli.Context) *relay.Relay {
	r := relay.New(newConfig(c))
	return r
}

func newBroker(c *cli.Context) broker.Broker {
	b := sns.New()
	return b
}
