package glint

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

	"github.com/mitchellh/go-glint/internal/flex"
)

// Document is the primary structure for managing and drawing components.
//
// A document represents a terminal window or session. The output can be set and
// components can be added, rendered, and drawn. All the methods on a Document
// are thread-safe unless otherwise documented. This allows you to draw,
// add components, replace components, etc. all while the render loop is active.
//
// Currently, this can only render directly to an io.Writer that expects to
// be a terminal session. In the future, we'll further abstract the concept
// of a "renderer" so that rendering can be done to other mediums as well.
type Document struct {
	mu          sync.Mutex
	w           io.Writer
	cols        uint
	rows        uint
	els         []Component
	refreshRate time.Duration
	lastHeight  uint
}

// SetOutput sets the location where rendering will be drawn.
//
// This should be a tty that supports ANSI escape sequences. In the future
// we'll better handle scenarios where ANSI escape sequences aren't supported.
func (d *Document) SetOutput(w io.Writer) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.w = w
}

// SetSize manually sets the size of the terminal window. If this is unset
// then we will automatically determine the terminal window size.
func (d *Document) SetSize(rows, cols uint) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.rows = rows
	d.cols = cols
}

// SetRefreshRate sets the rate at which output is rendered.
func (d *Document) SetRefreshRate(dur time.Duration) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.refreshRate = dur
}

// Append appends components to the document.
func (d *Document) Append(el ...Component) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.els = append(d.els, el...)
}

// Render starts a render loop that continues to render until the
// context is cancelled. This will render at the configured refresh rate.
// If the refresh rate is changed, it will not affect an active render loop.
// You must cancel and restart the render loop.
func (d *Document) Render(ctx context.Context) {
	d.mu.Lock()
	dur := d.refreshRate
	d.mu.Unlock()
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

	// Remove what we last drew. If what we last drew is greater than the number
	// of rows then we need to clear the screen.
	if d.lastHeight > 0 {
		if d.lastHeight <= rows {
			// Delete current line
			fmt.Fprint(d.w, b.Column(0).EraseLine(aec.EraseModes.All).ANSI)

			// Delete n lines above
			for i := uint(0); i < d.lastHeight-1; i++ {
				fmt.Fprint(d.w, b.Up(1).Column(0).EraseLine(aec.EraseModes.All).ANSI)
			}
		} else {
			fmt.Fprint(d.w, b.EraseDisplay(aec.EraseModes.All).EraseDisplay(aec.EraseMode(3)).Position(0, 0).ANSI)
		}
	}

	// Setup our root display which is our terminal. We don't set a height here
	// because we assume that the terminal has scrollback that the user can
	// use. If users want to lock into the terminal height they can use
	// a custom layout.
	config := flex.NewConfig()
	root := flex.NewNodeWithConfig(config)
	root.StyleSetWidth(float32(cols))

	// Build our render tree
	tree(root, Fragment(d.els...), rows, cols, false)

	// Calculate the layout
	flex.CalculateLayout(root, flex.Undefined, flex.Undefined, flex.DirectionLTR)

	// Render the tree
	renderTree(d.w, root, -1)

	// Store how much we drew
	height := uint(root.LayoutGetHeight())

	// If our component list is prefixed with finalized components, we
	// prune these out and do not re-render them.
	finalIdx := -1
	for i, el := range d.els {
		child := root.GetChild(i)
		if child == nil {
			break
		}

		// If the component is not finalized then we exit. If the
		// component doesn't match our expectations it means we hit
		// something weird and we exit too.
		ctx, ok := child.Context.(*parentContext)
		if !ok || ctx == nil || ctx.Component != el || !ctx.Finalized {
			break
		}

		// If this is finalized, then we have to subtract from the
		// height the height of this child since we're not going to redraw.
		// Then continue until we find one that isn't finalized.
		height -= uint(child.LayoutGetHeight())
		finalIdx = i
	}
	if finalIdx >= 0 {
		els := d.els[finalIdx+1:]
		d.els = make([]Component, len(els))
		copy(d.els, els)
	}

	// Store our last height which has now been processed for finalizations.
	d.lastHeight = height
}

var b = aec.EmptyBuilder
