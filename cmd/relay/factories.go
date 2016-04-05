package main

import (
	"errors"

	"github.com/lunohq/relay"
	"github.com/lunohq/relay/handler"
	"github.com/lunohq/relay/handler/sns"
	"github.com/lunohq/relay/handler/firehose"
	"github.com/lunohq/relay/handler/lambda"
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
		Handlers: newHandlers(c),
		TeamID: c.String("slack.team"),
		Token: c.String("slack.token"),
	}
}

func newRelay(c *cli.Context) *relay.Relay {
	r := relay.New(newConfig(c))
	return r
}

func newHandlers(c *cli.Context) (handlers []handler.Handler) {
	if t := c.String("sns.topic"); t != "" {
		handler := sns.New(sns.Options{
			TopicArn: t,
		})
		handlers = append(handlers, handler)
	} else if f := c.String("lambda.function"); f != "" {
		handler := lambda.New(lambda.Options{
			FunctionName: f,
			Qualifier: c.String("lambda.qualifier"),
		})
		handlers = append(handlers, handler)
	} else if d := c.String("firehose.stream"); d != "" {
		handler := firehose.New(firehose.Options{
			DeliveryStreamName: d,
		})
		handlers = append(handlers, handler)
	}

	if len(handlers) == 0 {
		must(errors.New("Must provide at least one handler config value"))
	}
	return handlers
}
