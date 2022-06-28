package zeromq

import (
	"github.com/goccy/go-json"
)

type Message struct {
	Topic    string
	Data     MessageData
	DataType string
}
type MessageHandler func(Message)
type MessageDataHandler func(data MessageData)
type MessageData string

func (d MessageData) ToString() string {
	var target string
	d.To(&target)
	return target
}

func (d MessageData) ToInt() int {
	var target int
	d.To(&target)
	return target
}

func (d MessageData) To(target interface{}) {
	json.Unmarshal([]byte(d), &target)
}

func NewMessage(zmqMsg []string) Message {
	topic := zmqMsg[0]
	msgLength := len(zmqMsg)

	var dataType string
	var data string

	if msgLength > 2 {
		dataType = zmqMsg[2]
		data = zmqMsg[1]
	} else {
		if msgLength > 1 {
			data = zmqMsg[1]
		} else {
			data = ""
		}
		dataType = "string"
	}

	return Message{
		Topic:    topic,
		Data:     MessageData(data),
		DataType: dataType,
	}

}
