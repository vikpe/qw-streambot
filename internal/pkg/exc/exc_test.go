package exc_test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/internal/pkg/exc"
)

func TestExecCommand(t *testing.T) {
	t.Run("invalid cmd", func(t *testing.T) {
		assert.Equal(t, "", exc.GetOutput("__invalid_cmd__"))
	})

	t.Run("valid cmd", func(t *testing.T) {
		nativeOutput, _ := exec.Command("ls").CombinedOutput()
		expect := strings.TrimSpace(string(nativeOutput))
		assert.Equal(t, expect, exc.GetOutput("ls"))
	})
}
