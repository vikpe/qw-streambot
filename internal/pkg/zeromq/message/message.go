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

type Handler func(msg Message)

type Message struct {
	Topic       string
	ContentType string
	Content     SerializedValue
}

func NewMessage(topic string, content any) Message {
	return Message{
		Topic:       topic,
		ContentType: fmt.Sprintf("%T", content),
		Content:     Serialize(content),
	}
}

func NewMessageFromFrames(frames []string) (Message, error) {
	frameCount := len(frames)
	const expectedFrameCount = 3 // topic, content type, content

	if frameCount != expectedFrameCount {
		err := errors.New(fmt.Sprintf("expected %d message frames, got %d", expectedFrameCount, frameCount))
		return Message{}, err
	}

	msg := Message{
		Topic:       frames[IndexTopic],
		ContentType: frames[IndexContentType],
		Content:     SerializedValue(frames[IndexContent]),
	}

	return msg, nil
}
