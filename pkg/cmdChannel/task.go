package cmdChannel

import (
	"fmt"
	"github.com/yunling101/ControllerManager/models"
	"github.com/yunling101/ControllerManager/models/asset"
	"github.com/yunling101/ControllerManager/models/event"
	"github.com/yunling101/ControllerManager/models/q"
	"github.com/yunling101/ControllerManager/models/task"
	"strings"
	"time"
)

const OnceTaskType = "once"

func (c *client) getTask(taskId string) (err error) {
	err = q.Table(c.task.TableName()).QueryOne(q.M{"uuid": taskId}, &c.task)
	if err != nil {
		err = fmt.Errorf("%s 任务不存在", taskId)
		return
	}
	err = models.DB.Model(&c.task).Association("Hosts").Find(&c.task.Hosts)
	if err != nil {
		err = fmt.Errorf("%s 任务主机查询出错", taskId)
		return
	}
	if len(c.task.Hosts) == 0 {
		err = fmt.Errorf("%s 任务绑定主机不能为空", taskId)
		return
	}
	for idx, h := range c.task.Hosts {
		_ = models.DB.Model(&h).Association("Certificate").Find(&c.task.Hosts[idx].Certificate)
	}
	return
}

func (c *client) taskDispatch(h asset.Host) {
	defer c.wg.Done()

	cli := new(client)
	cli.host = h
	if err := cli.newClient(); err != nil {
		c.fail = append(c.fail, cli.outputFail(err))
		return
	}

	r := cli.ssh.RunCommand(cli.cli, c.task.Scripts)
	if r.Error != nil {
		c.fail = append(c.fail, cli.outputFail(r.Error))
		return
	}

	c.success = append(c.success, cli.outputSuccess(string(r.Data.([]byte))))
}

func (c *client) once() {
	if c.task.TimerType == OnceTaskType {
		_ = q.Table(c.task.TableName()).UpdateOne(q.M{"uuid": c.task.Uuid}, q.M{"status": false})
	}
}

func (c *client) event() *client {
	if len(c.fail) != 0 {
		event.AddTask(c.task.Uuid, c.task.Name+" 执行失败", strings.Join(c.fail, ","))
	}
	return c
}

func (c *client) writeResult() *client {
	var taskResults task.Results
	result := append(append(c.success, c.fail...))
	count, _ := q.Table(taskResults.TableName()).Count(q.M{"uuid": c.task.Uuid})
	hosts := make([]string, 0)
	for _, h := range c.task.Hosts {
		hosts = append(hosts, h.Hostname)
	}
	taskResults = task.Results{
		Results:   strings.Join(result, ""),
		Hosts:     strings.Join(hosts, ","),
		Frequency: int(count) + 1,
		Status: fmt.Sprintf("Total: %v Success: %v Fail: %v",
			len(c.success)+len(c.fail), len(c.success), len(c.fail)),
		StartTime: c.start,
		EndTime:   time.Now(),
		Uuid:      c.task.Uuid,
	}
	_ = q.Table(taskResults.TableName()).InsertOne(&taskResults)
	return c
}
