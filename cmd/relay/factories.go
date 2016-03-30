package main

import (
	"github.com/lunohq/relay"
	"github.com/codegangsta/cli"
)

func newConfig(c *cli.Context) *relay.Config {
	return &relay.Config{TeamId: c.Args()[0], Token: c.Args()[1]}
}

func newRelay(c *cli.Context) *relay.Relay {
	r := relay.New(newConfig(c))
	return r
}
