package main

import (
	"context"
	"time"

	"github.com/mitchellh/go-glint"
	gc "github.com/mitchellh/go-glint/components"
)

func main() {
	// Create a progress bar and just render it periodically
	p := gc.Progress(100)

	// Create a text element that is updated with the time
	go func() {
		for {
			time.Sleep(50 * time.Millisecond)

			// Update our progress bar
			if p.Current() == p.Total() {
				p.SetCurrent(0)
			} else {
				p.Increment()
			}
		}
	}()

	d := glint.New()
	d.Append(
		glint.Text("Fixed to top"),
		p,
		glint.TextFunc(func(rows, cols uint) string {
			return time.Now().String()
		}),
		glint.Text("All with flexbox!"),
	)
	d.Render(context.Background())
}
