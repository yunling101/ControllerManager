package cmd

import (
	"github.com/urfave/cli/v2"
	"github.com/yunling101/ControllerManager/pkg/ctrManager"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func FlagsPrometheus() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "listen-address",
			Value: "0.0.0.0:9097",
			Usage: "Listening address and port",
		},
		&cli.StringFlag{
			Name:     "prometheus-rules-store-dir",
			Value:    "",
			Usage:    "Rules storage directory, rules directory in the prometheus configuration file",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "secret-key-file",
			Value: "conf/secret_key",
			Usage: "Secret key file for communication",
		},
		&cli.StringFlag{
			Name:  "prometheus-rules-suffix",
			Value: "rules",
			Usage: "Prometheus rules file suffix",
		},
		&cli.StringFlag{
			Name:  "prometheus-rules-prefix",
			Value: "yone",
			Usage: "Prometheus rules file prefix",
		},
		&cli.StringFlag{
			Name:  "prometheus-listen-address",
			Value: "http://127.0.0.1:9090",
			Usage: "The address that prometheus listens on",
		},
		&cli.StringFlag{
			Name:  "prometheus-basic-username",
			Value: "",
			Usage: "Prometheus basic authentication username",
		},
		&cli.StringFlag{
			Name:  "prometheus-basic-password",
			Value: "",
			Usage: "Prometheus basic authentication password",
		},
		&cli.BoolFlag{
			Name:  "prometheus-reload",
			Value: false,
			Usage: "Whether to reload prometheus after modification",
		},
	}
}

func ActionPrometheus() {
	go func() {
		if err := ctrManager.RunPrometheus(); err != nil {
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
