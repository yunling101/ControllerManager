package cmdChannel

import (
	"fmt"
	"github.com/yunling101/ControllerManager/models"
	"github.com/yunling101/ControllerManager/models/q"
)

func (c *client) getHost(host string) (err error) {
	err = q.Table(c.host.TableName()).QueryOne(q.M{"ip": host}, &c.host)
	if err != nil {
		err = fmt.Errorf("%s 主机不存在", host)
		return
	}
	err = models.DB.Model(&c.host).Association("Certificate").Find(&c.host.Certificate)
	if err != nil {
		err = fmt.Errorf("%s 主机凭证查询出错", host)
		return
	}
	return
}

func (c *client) newClient() (err error) {
	if len(c.host.Certificate) == 0 {
		err = fmt.Errorf("%s 主机绑定凭证不能为空", c.host.Ip)
		return
	}

	c.ssh = SSH{Host: c.host.Ip,
		Type:     c.host.Certificate[0].SshType,
		Password: c.host.Certificate[0].SshPassword,
		Port:     c.host.Certificate[0].SshPort,
		User:     c.host.Certificate[0].SshUser,
		KeyBody:  c.host.Certificate[0].SshKey,
	}
	c.cli, err = c.ssh.NewClient()

	return
}

func (c *client) commandRun(host, command string) {
	defer c.wg.Done()

	cli := new(client)
	cli.host.Hostname = host
	cli.host.Ip = host
	if err := cli.getHost(host); err != nil {
		cli.printOutputFail(err)
		return
	}

	if err := cli.newClient(); err != nil {
		cli.printOutputFail(err)
		return
	}

	response := cli.ssh.RunCommand(cli.cli, command)
	if response.Error != nil {
		cli.printOutputFail(response.Error)
		return
	}

	cli.printOutputSuccess(string(response.Data.([]byte)))
}
