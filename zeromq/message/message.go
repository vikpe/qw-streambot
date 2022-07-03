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
	Content     Content
	ContentType string
}

func NewMessage(topic string, content any) Message {
	return Message{
		Topic:       topic,
		Content:     NewContent(content),
		ContentType: fmt.Sprintf("%T", content),
	}
}

type Content string

func NewContent(content any) Content {
	return Content(Serialize(content))
}

func (d Content) ToString() string {
	var target string
	d.To(&target)
	return target
}

func (d Content) ToInt() int {
	var target int
	d.To(&target)
	return target
}

func (d Content) To(target interface{}) {
	Unserialize(string(d), &target)
}

func NewMessageFromFrames(frames []string) (Message, error) {
	frameCount := len(frames)

	switch frameCount {
	case 1:
		return NewMessage(frames[IndexTopic], ""), nil
	case 2:
		return NewMessage(frames[IndexTopic], frames[IndexContentType]), nil
	case 3:
		return Message{
			Topic:       frames[IndexTopic],
			Content:     NewContent(frames[IndexContent]),
			ContentType: frames[IndexContentType],
		}, nil
	}

	return Message{}, errors.New(fmt.Sprintf("expected 1-3 message frames, got %d", frameCount))
}
