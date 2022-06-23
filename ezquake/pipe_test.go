package ezquake_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/ezquake"
)

func pipePath(username string) string {
	return fmt.Sprintf("/tmp/ezquake_fifo_%s", username)
}
func resetPipe(username string) {
	os.Truncate(pipePath(username), 0)
}

func readPipe(username string) string {
	contentAsBytes, _ := os.ReadFile(pipePath(username))
	return string(contentAsBytes)
}

func TestPipeWriter_Write(t *testing.T) {
	username := "test"
	resetPipe(username)

	w := ezquake.NewWriter(username)

	w.Write("console;;")
	assert.Equal(t, "console;", readPipe(username))

	w.Write(" ")
	assert.Equal(t, "console;", readPipe(username))

	w.Write("lastscores")
	assert.Equal(t, "console;lastscores;", readPipe(username))
}
