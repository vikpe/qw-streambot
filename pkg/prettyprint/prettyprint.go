package prettyprint

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

const DELIMITER = "  "

var white = color.New(color.FgWhite).SprintFunc()

type PrettyPrinter struct {
	prefix string
}

func New(prefix string, colorCode color.Attribute) PrettyPrinter {
	return PrettyPrinter{
		prefix: color.New(colorCode).Sprint(prefix),
	}
}

func (pp PrettyPrinter) Println(args ...any) {
	pp.printTimestampAndPrefix()
	fmt.Println(args...)
}

func (pp PrettyPrinter) Printfln(format string, args ...any) {
	pp.Println(fmt.Sprintf(format, args...))
}

func (pp PrettyPrinter) Print(args ...any) {
	pp.printTimestampAndPrefix()
	fmt.Print(args...)
}

func (pp PrettyPrinter) Printf(format string, args ...any) {
	pp.Print(fmt.Sprintf(format, args...))
}

func (pp PrettyPrinter) printTimestampAndPrefix() {
	fmt.Print(timestamp(), DELIMITER, pp.prefix, DELIMITER)
}

func timestamp() string {
	currentTime := time.Now()
	format := "15:04:05"
	return white(currentTime.Format(format))
}
