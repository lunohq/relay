package main

import (
	"errors"

	"github.com/lunohq/relay"
	"github.com/lunohq/relay/handler"
	"github.com/lunohq/relay/handler/sns"
	"github.com/lunohq/relay/handler/firehose"
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
		Handler: newHandler(c),
		TeamID: c.String("slack.team"),
		Token: c.String("slack.token"),
	}
}

func newRelay(c *cli.Context) *relay.Relay {
	r := relay.New(newConfig(c))
	return r
}

func newHandler(c *cli.Context) handler.Handler {
	if t := c.String("sns.topic"); t != "" {
		return sns.New(sns.Options{
			TopicArn: t,
		})
	} else if d := c.String("firehose.stream"); d != "" {
		return firehose.New(firehose.Options{
			DeliveryStreamName: d,
		})
	}

	must(errors.New("Must provide at least one handler config value"))
	return nil
}
