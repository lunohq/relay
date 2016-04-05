package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "relay"
	app.Usage = "Relay messages from a slack bot to SNS"
	app.Action = mainAction
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "sns.topic",
			Value: "",
			Usage: "SNS Topic ARN to configure the SNS broker",
			EnvVar: "SNS_TOPIC",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func mainAction(c *cli.Context) {
	r := newRelay(c)
	r.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	sig := <-quit

	log.WithFields(log.Fields{
		"signal": sig,
	}).Info("terminating relay")

	if err := r.Shutdown(); err != nil {
		must(err)
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
