package main

import (
	"github.com/urfave/cli/v2"
	"github.com/yunling101/ControllerManager/common"
	"github.com/yunling101/ControllerManager/models"
	"github.com/yunling101/ControllerManager/pkg/cmd"
	"log"
	"os"
)

func run() (err error) {
	app := cli.NewApp()
	app.Name = "ctrManager"
	app.Usage = "A management toolset"
	app.Version = common.PrintVersion()
	app.CustomAppHelpTemplate = cmd.HelpTemplate()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "config",
			Value: "conf/config.yml",
			Usage: "Specify the configuration file location",
		},
	}
	app.Before = func(c *cli.Context) error {
		_ = common.LoadConfig(c.String("config"), "yaml")
		return models.Connect()
	}
	app.Flags = append(app.Flags, cmd.Flags()...)
	app.Action = func(context *cli.Context) error {
		cmd.Action(context)
		return nil
	}
	err = app.Run(os.Args)
	app.Commands = []*cli.Command{}
	return
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
