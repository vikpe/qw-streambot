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
	var topic string
	var dataType string
	var data string

	topic = rawMsg[0]

	if 3 == len(rawMsg) {
		dataType = rawMsg[2]
		data = rawMsg[1]
	} else {
		if 2 == len(rawMsg) {
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
