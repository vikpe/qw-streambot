package message

import "github.com/goccy/go-json"

type Message struct {
	Topic    string
	Data     Data
	DataType string
}
type Handler func(Message)
type DataHandler func(data Data)
type Data string

func (d Data) ToString() string {
	var target string
	d.To(&target)
	return target
}

func (d Data) ToInt() int {
	var target int
	d.To(&target)
	return target
}

func (d Data) To(target interface{}) {
	json.Unmarshal([]byte(d), &target)
}

func NewFromMultipart(zmqMsg []string) Message {
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
		Data:     Data(data),
		DataType: dataType,
	}
}
