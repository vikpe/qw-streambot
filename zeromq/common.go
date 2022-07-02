package zeromq

import (
	"time"
)

const TopicsAll = ""

func WaitForConnection() {
	time.Sleep(time.Millisecond * 10)
}
