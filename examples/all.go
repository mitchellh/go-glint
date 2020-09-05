package main

import (
	"context"
	"os"
	"time"

	"github.com/mitchellh/go-glint"
)

func main() {
	// Create a progress bar and just render it periodically
	p := glint.Progress(100)

	// Create a text element that is updated with the time
	timeEl := glint.Text("")
	go func() {
		for {
			time.Sleep(50 * time.Millisecond)

			// Update our progress bar
			if p.Current() == p.Total() {
				p.SetCurrent(0)
			} else {
				p.Increment()
			}

			// Update our time
			timeEl.Update(time.Now().String())
		}
	}()

	var d glint.Document
	d.SetOutput(os.Stdout)
	d.Add(
		glint.Text("Fixed to top"),
		p,
		timeEl,
		glint.Text("All with flexbox!"),
	)
	d.Render(context.Background())
}
