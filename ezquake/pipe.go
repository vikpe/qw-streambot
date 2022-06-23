package ezquake

import (
	"fmt"
	"os"
	"strings"
)

type Writer struct {
	username string
	queue    []string
}

func NewWriter(username string) Writer {
	return Writer{
		username: username,
		queue:    make([]string, 0),
	}
}

func (w Writer) Write(value string) error {
	strippedValue := strings.TrimSpace(value)

	if 0 == len(strippedValue) {
		return nil
	}

	terminatedValue := strings.TrimRight(strippedValue, ";") + ";"
	w.queue = append(w.queue, terminatedValue)
	return w.processQueue()
}

func (w Writer) processQueue() error {
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

func (w Writer) writeToPipe(value string) error {
	file, errOpen := os.OpenFile(w.path(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if errOpen != nil {
		return errOpen
	}

	defer file.Close()

	_, errWrite := file.WriteString(value)

	return errWrite
}

func (w Writer) path() string {
	return fmt.Sprintf("/tmp/ezquake_fifo_%s", w.username)
}
