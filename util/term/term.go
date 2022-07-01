package term

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

func NewPrettyPrinter(prefix string, colorCode color.Attribute) PrettyPrinter {
	return PrettyPrinter{
		prefix: color.New(colorCode).Sprint(prefix),
	}
}

func (pp PrettyPrinter) Println(args ...any) {
	fmt.Print(timestamp(), DELIMITER, pp.prefix, DELIMITER)
	fmt.Println(args...)
}

func (pp PrettyPrinter) Print(args ...any) {
	fmt.Print(timestamp(), DELIMITER, pp.prefix, DELIMITER)
	fmt.Print(args...)
}

func timestamp() string {
	currentTime := time.Now()
	format := "15:04:05"
	return white(currentTime.Format(format))
}
