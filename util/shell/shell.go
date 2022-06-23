package shell

import (
	"fmt"
	"os/exec"
	"strings"
)

func ExecCommand(command string) string {
	args := strings.Split(command, " ")
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("SHELL ERROR", err, cmd.String())
		return ""
	}

	return strings.TrimSpace(string(out))
}
