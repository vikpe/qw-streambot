package ezquake

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

type Controller struct {
	process Process
	pipe    PipeWriter
}

func NewController(process Process, pipe PipeWriter) Controller {
	return Controller{
		process: process,
		pipe:    pipe,
	}
}

func (c Controller) Command(cmd string) {
	if c.process.IsStarted() {
		c.pipe.Write(cmd)
	}
}

func (c Controller) Lastscores(duration time.Duration) {
	go func() {
		c.Command("toggleconsole")
		c.Command("lastscores")
		time.Sleep(duration)
		c.Command("toggleconsole")
	}()
}

func (c Controller) StaticText(text string) {
	trimmedText := strings.TrimSpace(fmt.Sprintf("bot_set_statictext %s", text))
	c.Command(trimmedText)

	if len(trimmedText) > 0 {
		textScale := staticTextScale(trimmedText)
		c.Command(fmt.Sprintf("hud_static_text_scale %f", textScale))
	}
}

func (c Controller) Showscores(duration time.Duration) {
	go func() {
		c.Command("+showscores")
		time.Sleep(duration)
		c.Command("-showscores")
	}()
}

func (c Controller) Play(playable string) {
	c.Command(fmt.Sprintf("qwurl %s", toQwurl(playable)))
}

func (c Controller) Weaponstats(duration time.Duration) {
	go func() {
		c.Command("bot_weaponstats_show")
		time.Sleep(duration)
		c.Command("bot_weaponstats_hide")
	}()
}

func staticTextScale(text string) float64 {
	scaleMin := 1.0
	scaleMax := 1.5
	lengthMin := float64(len("getquad semi")) // 12
	lengthMax := 3.0 * lengthMin

	lengthFactor := (float64(len(text)) - lengthMin) / (lengthMax - lengthMin)
	scale := scaleMax - (lengthFactor * (scaleMax - scaleMin))

	clamp := func(value float64, min_ float64, max_ float64) float64 {
		valueList := []float64{min_, max_, value}
		sort.Float64s(valueList)
		return valueList[1]
	}

	return Round(clamp(scaleMin, scale, scaleMax), 2)
}

func toQwurl(value string) string {
	if strings.Contains(value, "qw://") {
		return value
	} else if strings.Contains(value, "@") {
		return fmt.Sprintf("qw://%s/qtvplay", value)
	}

	return fmt.Sprintf("qw://%s/observe", value)
}

func Round(value float64, precision int) float64 {
	n := math.Pow(10, float64(precision))
	return math.Round(value*n) / n
}
