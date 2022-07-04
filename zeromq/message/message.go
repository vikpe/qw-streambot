package message

import (
	"errors"
	"fmt"
)

const (
	IndexTopic       = 0
	IndexContentType = 1
	IndexContent     = 2
)

type Handler func(Message)

type Message struct {
	Topic       string
	Content     SerializedValue
	ContentType string
}

func NewMessage(topic string, content any) Message {
	return Message{
		Topic:       topic,
		Content:     NewSerializedValue(content),
		ContentType: fmt.Sprintf("%T", content),
	}
}

type SerializedValue []byte

func NewSerializedValue(value any) SerializedValue {
	return Serialize(value)
}

func (d SerializedValue) ToString() string {
	var target string
	d.To(&target)
	return target
}

func (d SerializedValue) ToInt() int {
	var target int
	d.To(&target)
	return target
}

func (d SerializedValue) To(target interface{}) {
	Unserialize(d, target)
}

func NewMessageFromFrames(frames []string) (Message, error) {
	frameCount := len(frames)
	const expectedFrameCount = 3

	if frameCount != expectedFrameCount {
		return Message{}, errors.New(fmt.Sprintf("expected %d message frames, got %d", expectedFrameCount, frameCount))
	}

	msg := Message{
		Topic:       frames[IndexTopic],
		Content:     SerializedValue(frames[IndexContent]),
		ContentType: frames[IndexContentType],
	}

	fmt.Println("NewMessageFromFrames", "msg.Topic=", msg.Topic, msg.ContentType, "msg.Content=", msg.Content, "msg.Content.ToString()=", msg.Content.ToString())

	return msg, nil
}
