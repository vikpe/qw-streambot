package mock

import (
	"strings"

	"golang.org/x/exp/slices"
)

type ExecMock struct {
	Calls  []string
	Output map[string]string
}

func NewExecMock() ExecMock {
	return ExecMock{
		Calls:  make([]string, 0),
		Output: make(map[string]string, 0),
	}
}

func (m *ExecMock) Command(command string) string {
	m.Calls = append(m.Calls, command)
	args := strings.Split(command, " ")

	if response, ok := m.Output[args[0]]; ok {
		return response
	}
	return ""
}

func (m ExecMock) HasCommandCall(command string) bool {
	return slices.Contains(m.Calls, command)
}
