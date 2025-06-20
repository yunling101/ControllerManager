package event

import (
	"github.com/yunling101/ControllerManager/common"
	"github.com/yunling101/ControllerManager/models/q"
	"time"
)

type Event struct {
	Id           int       `json:"id"`
	EventId      string    `json:"event_id"`
	EventType    string    `json:"event_type"`
	EventTitle   string    `json:"event_title"`
	EventContent string    `json:"event_content"`
	Username     string    `json:"username"`
	Status       bool      `json:"status"`
	CreateTime   time.Time `json:"create_time"`
}

func (t *Event) TableName() string {
	return common.Config().Global.TableNamePrefix() + "_event"
}

func (t *Event) Add() {
	t.Status = true
	t.CreateTime = time.Now()
	_ = q.Table(t.TableName()).InsertOne(&t)
}

func (t *Event) SetEventType(eventType string) *Event {
	t.EventType = eventType
	return t
}

func (t *Event) SetEventId(eventId string) *Event {
	t.EventId = eventId
	return t
}

func (t *Event) SetEventTitle(eventTitle string) *Event {
	t.EventTitle = eventTitle
	return t
}

func (t *Event) SetEventContent(eventContent string) *Event {
	t.EventContent = eventContent
	return t
}

func (t *Event) SetUsername(username string) *Event {
	t.Username = username
	return t
}
