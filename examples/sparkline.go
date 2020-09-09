package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/mitchellh/go-glint"
	gc "github.com/mitchellh/go-glint/components"
)

func main() {
	max := 25
	min := 1

	values := make([]uint, 25)
	for i := range values {
		values[i] = uint(rand.Intn(max-min) + min)
	}

	// Create our sparkline
	sl := gc.Sparkline(values)
	sl.PeakStyle = []glint.StyleOption{glint.Color("green")}

	// Start up a timer that adds values
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			sl.Append(uint(rand.Intn(max-min) + min))
		}
	}()

	d := glint.New()
	d.Append(
		sl,
	)
	d.Render(context.Background())
}
