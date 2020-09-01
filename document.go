package dynamiccli

import (
	"bytes"
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
	cols        uint
	rows        uint
	els         []Element
	refreshRate time.Duration
	lastCount   uint
}

// SetOutput sets the location where rendering will be drawn.
func (d *Document) SetOutput(w io.Writer) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.w = w
}

func (d *Document) SetSize(rows, cols uint) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.rows = rows
	d.cols = cols
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
	cols := d.cols
	rows := d.rows
	if cols == 0 || rows == 0 {
		if f, ok := d.w.(*os.File); ok && sshterm.IsTerminal(int(f.Fd())) {
			ws, err := pty.GetsizeFull(f)
			if err == nil {
				rows = uint(ws.Rows)
				cols = uint(ws.Cols)
			}
		}
	}

	// We always have one less row than the size of the window because
	// we draw a newline at the end of every render.
	// NOTE(mitchellh): This is very fixable and we probably want to one day
	rows -= 1

	// We first render into a set of buffers so we can ensure that we only
	// render what we can see on the screen. If we render too many lines
	// then we'll create an infinite scrollback. This prevents that.
	var count uint
	var renders []*bytes.Buffer
	for i := len(d.els) - 1; i >= 0; i-- {
		el := d.els[i]

		var render bytes.Buffer
		thisCount := el.Render(&render, cols)
		nextCount := count + thisCount
		if nextCount > rows {
			break
		}

		count = nextCount
		renders = append(renders, &render)
	}

	// Clear the number of lines we rendered during the last pass. If this
	// is more than the rows that we have then we clear the rows.
	clear := d.lastCount
	if clear > rows {
		clear = rows
	}
	for i := uint(0); i < clear; i++ {
		fmt.Fprint(d.w, b.Up(1).Column(0).EraseLine(aec.EraseModes.All).ANSI)
	}

	// Go back and do our render
	for i := len(renders) - 1; i >= 0; i-- {
		io.Copy(d.w, bytes.NewReader(renders[i].Bytes()))
		fmt.Fprintln(d.w)
	}

	// Store how much we drew
	d.lastCount = count
}

var b = aec.EmptyBuilder
