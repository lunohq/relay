package main

import (
	"errors"

	"github.com/lunohq/relay"
	"github.com/lunohq/relay/broker"
	"github.com/lunohq/relay/broker/sns"
	"github.com/codegangsta/cli"
)

func newConfig(c *cli.Context) *relay.Config {
	if c.String("slack.token") == "" {
		must(errors.New("slack.token must be provided"))
		return nil
	}

	if c.String("slack.team") == "" {
		must(errors.New("slack.team must be provided"))
		return nil
	}

	return &relay.Config{
		Broker: newBroker(c),
		TeamID: c.String("slack.team"),
		Token: c.String("slack.token"),
	}
}

func newRelay(c *cli.Context) *relay.Relay {
	r := relay.New(newConfig(c))
	return r
}

func newBroker(c *cli.Context) broker.Broker {
	if t := c.String("sns.topic"); t != "" {
		return sns.New(sns.Options{
			TopicArn: t,
		})
	}
	must(errors.New("Must provide at least one broker config value"))
	return nil
}
