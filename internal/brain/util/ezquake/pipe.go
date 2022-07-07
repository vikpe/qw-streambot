package ezquake

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

type PipeWriter struct {
	username string
	mux      sync.Mutex
}

func NewPipeWriter(username string) *PipeWriter {
	return &PipeWriter{
		username: username,
		mux:      sync.Mutex{},
	}
}

func (w *PipeWriter) Write(value string) error {
	trimmedValue := strings.TrimSpace(value)

	if 0 == len(trimmedValue) {
		return nil
	}

	terminatedValue := strings.TrimRight(trimmedValue, ";") + ";"
	return w.writeToPipe(terminatedValue)
}

func (w *PipeWriter) writeToPipe(value string) error {
	w.mux.Lock()
	defer w.mux.Unlock()

	file, errOpen := os.OpenFile(w.path(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if errOpen != nil {
		return errOpen
	}

	defer file.Close()

	_, errWrite := file.WriteString(value)

	return errWrite
}

func (w *PipeWriter) path() string {
	return fmt.Sprintf("/tmp/ezquake_fifo_%s", w.username)
}
