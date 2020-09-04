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

	"github.com/mitchellh/go-dynamic-cli/internal/flex"
)

// TODO: docs
type Document struct {
	mu          sync.Mutex
	w           io.Writer
	cols        uint
	rows        uint
	els         []Component
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

func (d *Document) Add(el ...Component) {
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

	// Remove what we last drew. If what we last drew is greater than the number
	// of rows then we need to clear the screen.
	if d.lastCount > 0 {
		if d.lastCount <= rows {
			// Delete current line
			fmt.Fprint(d.w, b.Column(0).EraseLine(aec.EraseModes.All).ANSI)

			// Delete n lines above
			for i := uint(0); i < d.lastCount-1; i++ {
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
	tree(root, Fragment(d.els...), rows, cols)

	// Calculate the layout
	flex.CalculateLayout(root, flex.Undefined, flex.Undefined, flex.DirectionLTR)

	// Debug. Flip this to true to see flexbox calculations.
	if false {
		fmt.Printf("rows: %d\n", rows)
		fmt.Printf("cols: %d\n", cols)
		fmt.Printf("root left: %f\n", root.LayoutGetLeft())     // 0
		fmt.Printf("root top: %f\n", root.LayoutGetTop())       // 0
		fmt.Printf("root width: %f\n", root.LayoutGetWidth())   // 200
		fmt.Printf("root height: %f\n", root.LayoutGetHeight()) // 200
		for i := 0; ; i++ {
			child := root.GetChild(i)
			if child == nil {
				break
			}

			fmt.Printf("child %d left: %f\n", i, child.LayoutGetLeft())     // 0
			fmt.Printf("child %d top: %f\n", i, child.LayoutGetTop())       // 0
			fmt.Printf("child %d width: %f\n", i, child.LayoutGetWidth())   // 200
			fmt.Printf("child %d height: %f\n", i, child.LayoutGetHeight()) // 200
		}
		os.Exit(1)
	}

	// Render the tree
	renderTree(d.w, root, -1)

	// Store how much we drew
	d.lastCount = uint(root.LayoutGetHeight())
}

var b = aec.EmptyBuilder
