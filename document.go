package dynamiccli

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/creack/pty"
	"github.com/morikuni/aec"
	sshterm "golang.org/x/crypto/ssh/terminal"
)

// TODO: docs
type Document struct {
	mu          sync.Mutex
	w           io.Writer
	width       uint
	els         []Element
	lineCount   uint
	refreshRate time.Duration
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

func (d *Document) SetRefreshRate(dur time.Duration) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.refreshRate = dur
}

func (d *Document) Add(el ...Element) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.els = append(d.els, el...)
}

// Render starts a render loop that continues to render until the
// context is cancelled.
func (d *Document) Render(ctx context.Context) {
	dur := d.refreshRate
	if dur == 0 {
		dur = time.Second / 12
	}

	t := time.NewTicker(dur)
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
//
// If a manual size is not configured, this will recalcualte the window
// size on each call. This typically requires a syscall. This is a bit
// expensive but something we can optimize in the future if it ends up being
// a real source of FPS issues.
func (d *Document) RenderFrame() {
	d.mu.Lock()
	defer d.mu.Unlock()

	// If we don't have a writer set, then don't render anything.
	if d.w == nil {
		return
	}

	// Detect if we had a window size change
	width := d.width
	if width == 0 {
		if f, ok := d.w.(*os.File); ok && sshterm.IsTerminal(int(f.Fd())) {
			ws, err := pty.GetsizeFull(f)
			if err == nil {
				width = uint(ws.Cols)
			}
		}
	}

	// Delete prior output.
	// TODO(mitchellh): on resizing to a smaller width, terminals will
	// typically word wrap. We don't currently detect cursor location to clear
	// this.
	if d.lineCount > 0 {
		fmt.Fprint(d.w, b.Up(d.lineCount).Column(0).EraseLine(aec.EraseModes.All).ANSI)
	}

	// Reset our line count to zero and start rerendering
	d.lineCount = 0

	// Go through each element and output.
	for _, el := range d.els {
		d.lineCount += el.Render(d.w, width)
		fmt.Fprint(d.w, "\n")
	}
}

var b = aec.EmptyBuilder
