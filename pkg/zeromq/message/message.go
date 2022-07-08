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
		Content:     Serialize(content),
		ContentType: fmt.Sprintf("%T", content),
	}
}

func NewMessageFromFrames(frames []string) (Message, error) {
	frameCount := len(frames)
	const expectedFrameCount = 3 // topic, data type, data

	if frameCount != expectedFrameCount {
		return Message{}, errors.New(fmt.Sprintf("expected %d message frames, got %d", expectedFrameCount, frameCount))
	}

	msg := Message{
		Topic:       frames[IndexTopic],
		Content:     SerializedValue(frames[IndexContent]),
		ContentType: frames[IndexContentType],
	}

	return msg, nil
}
