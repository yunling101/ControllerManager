package cmd

import (
	"github.com/urfave/cli/v2"
	"github.com/yunling101/ControllerManager/models"
	"github.com/yunling101/ControllerManager/pkg/ctrManager"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func HelpTemplate() string {
	return `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{if .Version}}{{if not .HideVersion}}

VERSION:
   {{.Version}}{{end}}{{end}}{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}{{if len .Authors}}

AUTHOR:
   {{range .Authors}}{{ . }}{{end}}{{end}}{{if .VisibleCommands}}
{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`
}

func Flags() (flag []cli.Flag) {
	flag = append(flag,
		&cli.StringFlag{
			Name:  "listen",
			Value: "0.0.0.0:9096",
			Usage: "Listening address",
		},
	)
	return
}

func Action(c *cli.Context) {
	go newCronTick()
	go func() {
		if err := ctrManager.Run(c.String("listen")); err != nil {
			log.Fatalln(err)
		}
	}()

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sign
		_ = models.SqlDB.Close()
		os.Exit(0)
	}()

	select {}
}