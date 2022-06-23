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

func (pp PrettyPrinter) Print(args ...any) {
	fmt.Print(pp.prefix, DELIMITER, timestamp(), DELIMITER)
	fmt.Println(args...)
}

func timestamp() string {
	currentTime := time.Now()
	format := "15:04:05"
	return white(currentTime.Format(format))
}
