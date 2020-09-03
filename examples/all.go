package main

import (
	"context"
	"os"
	"time"

	"github.com/mitchellh/go-dynamic-cli"
)

func main() {
	// Create a progress bar and just render it periodically
	p := dynamiccli.Progress(100)

	// Create a text element that is updated with the time
	timeEl := dynamiccli.Text("")
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

	var d dynamiccli.Document
	d.SetOutput(os.Stdout)
	d.Add(
		dynamiccli.Text("Fixed to top"),
		p,
		timeEl,
		dynamiccli.Text("All with flexbox!"),
	)
	d.Render(context.Background())
}
