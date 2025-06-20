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
	app.Name = "prometheus-plug"
	app.Usage = "A prometheus plugin tools"
	app.CustomAppHelpTemplate = cmd.HelpTemplate()
	app.Flags = cmd.FlagsPrometheus()
	app.Before = func(c *cli.Context) error {
		return common.LoadFlagsPrometheus(c)
	}
	app.Action = func(c *cli.Context) error {
		cmd.ActionPrometheus()
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
