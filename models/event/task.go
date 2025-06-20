package event

const taskName = "任务"

func AddTask(uuid, title, content string) {
	event := new(Event)
	event.SetEventId(uuid).
		SetEventType(taskName).
		SetUsername("task").
		SetEventTitle("【" + taskName + "】" + title).
		SetEventContent(content).
		Add()
}

func AddTaskError(uuid, title string, err error) error {
	if err != nil {
		AddTask(uuid, title, err.Error())
	}
	return err
}
