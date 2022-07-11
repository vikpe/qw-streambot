package zeromq

import (
	"time"
)

const TopicsAll = ""
const ConnectionGraceTimeout = time.Millisecond * 20

func WaitForConnection() {
	time.Sleep(ConnectionGraceTimeout)
}
