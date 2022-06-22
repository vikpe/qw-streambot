package ezquake

import (
	"fmt"
	"os"
	"strings"
)

type PipeWriter struct {
	username string
	queue    []string
}

func NewPipeWriter(username string) PipeWriter {
	return PipeWriter{
		username: username,
		queue:    make([]string, 0),
	}
}

func (w PipeWriter) Write(value string) error {
	strippedValue := strings.TrimSpace(value)

	if 0 == len(strippedValue) {
		return nil
	}

	terminatedValue := strings.TrimRight(strippedValue, ";") + ";"
	w.queue = append(w.queue, terminatedValue)
	return w.processQueue()
}

func (w PipeWriter) processQueue() error {
	for {
		if 0 == len(w.queue) {
			break
		}

		value := w.queue[0]
		err := w.writeToPipe(value)

		if err != nil {
			return err
		}

		if 1 == len(w.queue) {
			break
		}

		w.queue = w.queue[1:]
	}

	return nil
}

func (w PipeWriter) writeToPipe(value string) error {
	file, errOpen := os.OpenFile(w.path(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if errOpen != nil {
		return errOpen
	}

	defer file.Close()

	_, errWrite := file.WriteString(value)

	return errWrite
}

func (w PipeWriter) path() string {
	return fmt.Sprintf("/tmp/ezquake_fifo_%s", w.username)
}
