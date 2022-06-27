package zeromq

import (
	"github.com/goccy/go-json"
)

type Message struct {
	Topic    string
	Data     MessageJsonData
	DataType string
}

type MessageHandler func(Message)

type MessageJsonData string

func (d MessageJsonData) ToString() string {
	var target string
	d.To(&target)
	return target
}

func (d MessageJsonData) ToInt() int {
	var target int
	d.To(&target)
	return target
}

func (d MessageJsonData) To(target interface{}) {
	json.Unmarshal([]byte(d), &target)
}

func ParseMessage(rawMsg []string) Message {
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

	return Message{
		Topic:    topic,
		Data:     MessageJsonData(data),
		DataType: dataType,
	}

}
