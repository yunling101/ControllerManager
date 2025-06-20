package cmdChannel

import "github.com/urfave/cli/v2"

func Commands() (cmd []*cli.Command) {
	cmd = append(cmd, &cli.Command{
		Name:  "batch",
		Usage: "Batch command",
		Action: func(c *cli.Context) (err error) {
			err = New().CmdRun(c.String("command"),
				c.String("hosts"))
			return
		},
		Flags: []cli.Flag{
			StringFlag("command", "", "Command to be executed"),
			StringFlag("hosts", "", "List of hosts for batch execution"),
		},
	})
	cmd = append(cmd, &cli.Command{
		Name:  "task",
		Usage: "Scheduled task",
		Action: func(c *cli.Context) (err error) {
			err = New().TaskRun(c.String("task_id"))
			return
		},
		Flags: []cli.Flag{
			StringFlag("task_id", "", "Task ID to be executed"),
		},
	})
	return
}

func StringFlag(name, value, usage string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:     name,
		Value:    value,
		Usage:    usage,
		Required: true,
	}
}
