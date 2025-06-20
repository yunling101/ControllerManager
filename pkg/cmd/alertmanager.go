package cmd

import (
	"github.com/urfave/cli/v2"
	"github.com/yunling101/ControllerManager/pkg/ctrManager"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func FlagsAlertmanager() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "listen-address",
			Value: "0.0.0.0:9098",
			Usage: "Listening address and port",
		},
		&cli.StringFlag{
			Name:     "alertmanager-base-dir",
			Value:    "",
			Usage:    "Alertmanager installation directory (configuration file)",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "secret-key-file",
			Value: "conf/secret_key",
			Usage: "Secret key file for communication",
		},
		&cli.StringFlag{
			Name:  "alertmanager-config-name",
			Value: "alertmanager.yml",
			Usage: "Alertmanager configuration file name",
		},
		&cli.StringFlag{
			Name:  "alertmanager-listen-address",
			Value: "http://127.0.0.1:9093",
			Usage: "The address that alertmanager listens on",
		},
		&cli.StringFlag{
			Name:  "alertmanager-basic-username",
			Value: "",
			Usage: "Alertmanager basic authentication username",
		},
		&cli.StringFlag{
			Name:  "alertmanager-basic-password",
			Value: "",
			Usage: "Alertmanager basic authentication password",
		},
		&cli.BoolFlag{
			Name:  "alertmanager-reload",
			Value: false,
			Usage: "Whether to reload alertmanager after modification",
		},
	}
}

func ActionAlertmanager() {
	go func() {
		if err := ctrManager.RunAlertmanager(); err != nil {
			log.Fatalln(err)
		}
	}()

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sign
		os.Exit(0)
	}()

	select {}
}
