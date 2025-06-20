package cmdChannel

import (
	"fmt"
	"github.com/yunling101/ControllerManager/models/asset"
	"github.com/yunling101/ControllerManager/models/event"
	"github.com/yunling101/ControllerManager/models/task"
	"golang.org/x/crypto/ssh"
	"strings"
	"sync"
	"time"
)

type client struct {
	wg      sync.WaitGroup
	start   time.Time
	host    asset.Host
	task    task.Task
	cli     *ssh.Client
	ssh     SSH
	success []string
	fail    []string
}

func New() *client {
	return &client{start: time.Now()}
}

func (c *client) CmdRun(command, hosts string) error {
	for _, h := range strings.Split(hosts, ",") {
		c.wg.Add(1)
		go c.commandRun(h, command)
	}
	c.wg.Wait()

	return nil
}

func (c *client) TaskRun(taskId string) (err error) {
	err = event.AddTaskError(taskId, "获取失败", c.getTask(taskId))
	if err != nil {
		return
	}

	for _, h := range c.task.Hosts {
		c.wg.Add(1)
		go c.taskDispatch(*h)
	}
	c.wg.Wait()

	c.writeResult().once()
	return
}

func (c *client) outputFail(err error) string {
	return fmt.Sprintf("OUTPUT ==> %s:[fail]\n%s\n\n", c.host.Hostname, err.Error())
}

func (c *client) outputSuccess(result string) string {
	return fmt.Sprintf("OUTPUT ==> %s:[success]\n%s\n", c.host.Hostname, result)
}

func (c *client) printOutputFail(err error) {
	fmt.Printf("\033[1m\033[1;31m\nOUTPUT ==> %s\033[1m\033[0m\n%s\n\n", c.host.Hostname, err.Error())
}

func (c *client) printOutputSuccess(result string) {
	fmt.Printf("\033[1m\033[1;32m\nOUTPUT ==> %s\033[1m\033[0m\n%s\n", c.host.Hostname, result)
}
