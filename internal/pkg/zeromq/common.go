package zeromq

import (
	"time"
)

const TopicsAll = ""
const ConnectionGraceTimeout = time.Millisecond * 10

func WaitForConnection() {
	time.Sleep(ConnectionGraceTimeout)
}
