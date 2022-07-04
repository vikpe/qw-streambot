package term_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/util/term"
)

func TestPrettyPrinter_Println(t *testing.T) {
	testRun := func() {
		printer := term.NewPrettyPrinter("alpha", color.FgCyan)
		printer.Println("hello", 123)
	}

	expect := fmt.Sprintf("%s  alpha  hello 123\n", time.Now().Format("15:04:05"))
	output := getFuncStdOutput(testRun)
	assert.Equal(t, expect, output)
}

func TestPrettyPrinter_Print(t *testing.T) {
	testRun := func() {
		printer := term.NewPrettyPrinter("alpha", color.FgCyan)
		printer.Print("hello", 123)
	}

	expect := fmt.Sprintf("%s  alpha  hello123", time.Now().Format("15:04:05"))
	output := getFuncStdOutput(testRun)
	assert.Equal(t, expect, output)
}

func getFuncStdOutput(f func()) string {
	rescueStderr := os.Stderr
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	f()

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = rescueStderr
	os.Stdout = rescueStdout

	return string(out)
}
