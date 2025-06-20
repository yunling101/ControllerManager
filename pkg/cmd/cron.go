package cmd

import (
	"github.com/robfig/cron/v3"
	"github.com/yunling101/ControllerManager/logf"
	"github.com/yunling101/ControllerManager/models/event"
	"github.com/yunling101/ControllerManager/models/q"
	"github.com/yunling101/ControllerManager/models/task"
	"github.com/yunling101/ControllerManager/pkg/cache"
	"github.com/yunling101/ControllerManager/pkg/cmdChannel"
	"time"
)

const (
	loopInterval       = 90 // seconds
	cleanIntervalTimes = 80
)

type tick struct {
	counter   int
	task      task.Task
	taskQueue chan *task.Task
	Cron      *cron.Cron
	log       *logf.Log
}

type OnceJob struct {
	tick     *tick
	entryID  cron.EntryID
	executed bool
}

func (j *OnceJob) Run() {
	if !j.executed {
		j.tick.log.ErrorIf(cmdChannel.New().TaskRun(j.tick.task.Uuid))
		j.executed = true
	} else {
		if int(j.entryID) != 0 {
			j.tick.Cron.Remove(j.entryID)
			cache.Cache.Delete(j.tick.task.Uuid)
		}
	}
}

func newCronTick() {
	queue := make(chan *task.Task)

	// cron.WithSeconds()
	t := tick{taskQueue: queue, Cron: cron.New(), log: logf.Logger(logf.Parser(logf.TASK))}
	t.loopTasks()

	t.Cron.Start()
	t.receiveTask()
}

func (t *tick) getTasks() (tasks []task.Task) {
	_ = q.Table(t.task.TableName()).QueryMany(q.M{"status": true}, &tasks)
	return
}

func (t *tick) loopTasks() {
	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(loopInterval)).C
		for {
			<-ticker

			tasks := t.getTasks()
			if len(tasks) != 0 {
				if t.counter >= cleanIntervalTimes {
					go t.loopClearQueue(tasks)
					t.counter = 0
				}
				for _, h := range tasks {
					go func(task task.Task) {
						t.taskQueue <- &task
					}(h)
				}
				t.counter++
			}
		}
	}()
}

func (t *tick) loopClearQueue(tasks []task.Task) {
	for key, id := range cache.Cache.List() {
		isExist := false
		for _, h := range tasks {
			if key == h.Uuid {
				isExist = true
				break
			}
		}
		if !isExist {
			t.Cron.Remove(cron.EntryID(id))
			cache.Cache.Delete(key)
		}
	}
}

func (t *tick) receiveTask() {
	for {
		itemTask := <-t.taskQueue

		var err error
		if itemTask.TimerType == cmdChannel.OnceTaskType {
			err = t.processOnceTask(itemTask)
		} else {
			err = t.processTask(itemTask)
		}
		t.log.ErrorIf(event.AddTaskError(itemTask.Uuid, itemTask.Name+" 调度失败", err))
	}
}

func (t *tick) processOnceTask(itemTask *task.Task) error {
	if !cache.Cache.IsExist(itemTask.Uuid) {
		onceTaskProcess := &OnceJob{tick: &tick{task: *itemTask, Cron: t.Cron, log: t.log}}
		id, err := t.Cron.AddJob(itemTask.ExecutionTime, onceTaskProcess)
		if err != nil {
			return err
		}
		onceTaskProcess.entryID = id
		cache.Cache.Set(itemTask.Uuid, int(id))
	}
	return nil
}

func (t *tick) processTask(itemTask *task.Task) error {
	if !cache.Cache.IsExist(itemTask.Uuid) {
		id, err := t.Cron.AddFunc(itemTask.ExecutionTime, func() {
			t.log.ErrorIf(cmdChannel.New().TaskRun(itemTask.Uuid))
		})
		if err != nil {
			return err
		}
		cache.Cache.Set(itemTask.Uuid, int(id))
	}
	return nil
}
