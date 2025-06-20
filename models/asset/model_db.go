package asset

import "github.com/yunling101/ControllerManager/common"

type Certificate struct {
	Id          int    `json:"id"`
	SshType     string `json:"ssh_type"`
	SshPort     int    `json:"ssh_port"`
	SshUser     string `json:"ssh_user"`
	SshPassword string `json:"ssh_password"`
	SshKey      string `json:"ssh_key"`
}

type Host struct {
	Id          int            `json:"id"`
	Hostname    string         `json:"hostname"`
	Ip          string         `json:"ip"`
	Certificate []*Certificate `json:"certificate" gorm:"many2many:yone_hosts_certificate;JoinForeignKey:hosts_id;joinReferences:certificate_id"`
}

func (h *Host) TableName() string {
	return common.Config().Global.TableNamePrefix() + "_hosts"
}

func (h *Certificate) TableName() string {
	return common.Config().Global.TableNamePrefix() + "_certificate"
}
