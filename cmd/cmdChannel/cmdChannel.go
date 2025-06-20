package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/yunling101/ControllerManager/common"
	"github.com/yunling101/ControllerManager/models"
	"github.com/yunling101/ControllerManager/pkg/cmdChannel"
	"os"
)

// https://github.com/golang/glog
// https://github.com/kubernetes/klog
// https://github.com/bigkevmcd/go-configparser
// https://github.com/go-ini/ini
// https://github.com/BurntSushi/toml

func run() (err error) {
	app := cli.NewApp()
	app.Name = "cmdChannel"
	app.Usage = "A batch command tools"
	app.Before = func(c *cli.Context) (err error) {
		_ = common.LoadConfig(c.String("config"), "ini")
		err = models.Connect()
		return
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "config",
			Value: "webserver/config/default.ini",
			Usage: "Specify the configuration file location",
		},
	}
	app.Commands = cmdChannel.Commands()
	app.After = func(c *cli.Context) (err error) {
		return models.SqlDB.Close()
	}
	err = app.Run(os.Args)
	return
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(0)
	}
}
