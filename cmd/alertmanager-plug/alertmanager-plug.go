package main

import (
	"github.com/urfave/cli/v2"
	"github.com/yunling101/ControllerManager/common"
	"github.com/yunling101/ControllerManager/pkg/cmd"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "alertmanager-plug"
	app.Usage = "A alertmanager plugin tools"
	app.CustomAppHelpTemplate = cmd.HelpTemplate()
	app.Flags = cmd.FlagsAlertmanager()
	app.Before = func(c *cli.Context) error {
		return common.LoadFlagsAlertmanager(c)
	}
	app.Action = func(c *cli.Context) error {
		cmd.ActionAlertmanager()
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
