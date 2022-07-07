package mock

type PublisherMock struct {
	SendMessageCalls [][]any
}

func NewPublisherMock() PublisherMock {
	return PublisherMock{}
}

func (s *PublisherMock) SendMessage(topic string, args ...any) {
	s.SendMessageCalls = append(s.SendMessageCalls, append([]any{topic}, args...))
}
