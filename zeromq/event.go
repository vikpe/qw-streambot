package zeromq

import (
	"github.com/goccy/go-json"
)

type Event struct {
	Topic    string
	Data     EventData
	DataType string
}

type EventHandler func(Event)
type EventDataHandler func(data EventData)

type EventData string

func (d EventData) ToString() string {
	var target string
	d.To(&target)
	return target
}

func (d EventData) ToInt() int {
	var target int
	d.To(&target)
	return target
}

func (d EventData) To(target interface{}) {
	json.Unmarshal([]byte(d), &target)
}

func ParseEvent(rawMsg []string) Event {
	topic := rawMsg[0]
	msgPartCount := len(rawMsg)

	var dataType string
	var data string

	if msgPartCount > 2 {
		dataType = rawMsg[2]
		data = rawMsg[1]
	} else {
		if msgPartCount > 1 {
			data = rawMsg[1]
		} else {
			data = ""
		}
		dataType = "string"
	}

	return Event{
		Topic:    topic,
		Data:     EventData(data),
		DataType: dataType,
	}

}
