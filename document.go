package dynamiccli

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
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

	// Remove what we last drew. If what we last drew is greater than the number
	// of rows then we need to clear the screen.
	if d.lastCount <= rows {
		for i := uint(0); i < d.lastCount; i++ {
			fmt.Fprint(d.w, b.Up(1).Column(0).EraseLine(aec.EraseModes.All).ANSI)
		}
	} else {
		// TODO: clear display
	}

	// Render our elements
	var b strings.Builder
	for _, el := range d.els {
		b.WriteString(el.Render(cols))
		b.WriteRune('\n')
	}

	// Store how much we drew
	d.lastCount = uint(countLines(b.String()))

	// Draw
	io.Copy(d.w, strings.NewReader(b.String()))
}

func countLines(s string) int {
	count := strings.Count(s, "\n")

	// If the last character isn't a newline, we have to add one since we'll
	// always have one more line than newline characters.
	if len(s) > 0 && s[len(s)-1] != '\n' {
		count++
	}
	return count
}

var b = aec.EmptyBuilder
