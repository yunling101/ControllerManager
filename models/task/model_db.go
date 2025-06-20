package task

import (
	"github.com/yunling101/ControllerManager/common"
	"github.com/yunling101/ControllerManager/models/asset"
	"time"
)

type Task struct {
	Id            int           `json:"id"`
	Name          string        `json:"name"`
	Uuid          string        `json:"uuid"`
	TimerType     string        `json:"timer_type"`
	Scripts       string        `json:"scripts"`
	Status        bool          `json:"status"`
	Hosts         []*asset.Host `json:"hosts" gorm:"many2many:yone_task_timer_hosts;JoinForeignKey:tasktimer_id;joinReferences:hosts_id"`
	Notify        string        `json:"notify"`
	RetryCount    int           `json:"retry_count"`
	ExecutionTime string        `json:"execution_time"`
}

type Results struct {
	Id        int       `json:"id"`
	Uuid      string    `json:"uuid"`
	Results   string    `json:"results"`
	Hosts     string    `json:"hosts"`
	Status    string    `json:"status"`
	Frequency int       `json:"frequency"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func (t *Task) TableName() string {
	return common.Config().Global.TableNamePrefix() + "_task_timer"
}

func (t *Results) TableName() string {
	return common.Config().Global.TableNamePrefix() + "_task_timer_results"
}
