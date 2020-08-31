package dynamiccli

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/morikuni/aec"
)

// TODO: docs
type Document struct {
	mu        sync.Mutex
	w         io.Writer
	width     uint
	els       []Element
	lineCount uint
}

// SetOutput sets the location where rendering will be drawn.
func (d *Document) SetOutput(w io.Writer) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.w = w
}

func (d *Document) SetWidth(w uint) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.width = w
}

func (d *Document) Add(el ...Element) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.els = append(d.els, el...)
}

// Render starts a render loop that continues to render until the
// context is cancelled.
func (d *Document) Render(ctx context.Context) {
	t := time.NewTicker(time.Second / 6)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-t.C:
			d.RenderFrame()
		}
	}
}

// RenderFrame will render a single frame and return.
func (d *Document) RenderFrame() {
	d.mu.Lock()
	defer d.mu.Unlock()

	// If we don't have a writer set, then don't render anything.
	if d.w == nil {
		return
	}

	if d.lineCount > 0 {
		// Delete prior output
		fmt.Fprint(d.w, b.Up(d.lineCount).Column(0).EraseLine(aec.EraseModes.All).ANSI)
	}

	// Reset our line count to zero and start rerendering
	d.lineCount = 0

	// Go through each element and output.
	for _, el := range d.els {
		d.lineCount += el.Render(d.w, d.width)
		fmt.Fprint(d.w, "\n")
	}
}

var b = aec.EmptyBuilder
