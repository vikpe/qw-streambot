package zeromq

type Publisher struct {
	address string
}

func NewPublisher(address string) Publisher {
	return Publisher{address: address}
}

func (p Publisher) SendMessage(topic string, data any) {
	PubSendMessage(p.address, topic, data)
}
