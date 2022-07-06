package zeromq

import (
	"time"
)

const TopicsAll = ""

type EventHandler = func(topic string, data ...any)

func WaitForConnection() {
	time.Sleep(time.Millisecond * 10)
}
