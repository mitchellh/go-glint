package components

import (
	"time"

	"github.com/mitchellh/go-dynamic-cli"
)

// Stopwatch creates a new stopwatch component that starts at the given time.
func Stopwatch(start time.Time) *StopwatchComponent {
	return &StopwatchComponent{
		start: start,
	}
}

type StopwatchComponent struct {
	start time.Time
}

func (c *StopwatchComponent) Body() dynamiccli.Component {
	return dynamiccli.Text(time.Now().Sub(c.start).Truncate(100 * time.Millisecond).String())
}
