package prettyfmt

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

const DELIMITER = "  "

var white = color.New(color.FgWhite).SprintFunc()

type PrettyFmt struct {
	prefix string
}

func New(prefix string, colorCode color.Attribute) PrettyFmt {
	return PrettyFmt{
		prefix: color.New(colorCode).Sprint(prefix),
	}
}

func (pp PrettyFmt) Println(args ...any) {
	pp.printTimestampAndPrefix()
	fmt.Println(args...)
}

func (pp PrettyFmt) Printfln(format string, args ...any) {
	pp.Println(fmt.Sprintf(format, args...))
}

func (pp PrettyFmt) Print(args ...any) {
	pp.printTimestampAndPrefix()
	fmt.Print(args...)
}

func (pp PrettyFmt) Printf(format string, args ...any) {
	pp.Print(fmt.Sprintf(format, args...))
}

func (pp PrettyFmt) printTimestampAndPrefix() {
	fmt.Print(timestamp(), DELIMITER, pp.prefix, DELIMITER)
}

func timestamp() string {
	currentTime := time.Now()
	format := "15:04:05"
	return white(currentTime.Format(format))
}
