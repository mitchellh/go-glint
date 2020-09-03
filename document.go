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
		fmt.Fprint(d.w, b.EraseDisplay(aec.EraseModes.All).EraseDisplay(aec.EraseMode(3)).Position(0, 0).ANSI)
	}

	// Setup our root display which is our terminal
	config := flex.NewConfig()
	root := flex.NewNodeWithConfig(config)
	root.StyleSetWidth(float32(cols))
	root.StyleSetMaxHeight(float32(rows))
	//root.StyleSetMaxHeight(5)
	root.StyleSetOverflow(flex.OverflowHidden)

	// Render our elements
	elCache := make([]*measureContext, len(d.els))
	for idx, el := range d.els {
		// If the element wants the terminal size, give it.
		if el, ok := el.(ElementTerminalSizer); ok {
			el.SetTerminalSize(rows, cols)
		}

		// Setup our node
		node := flex.NewNodeWithConfig(config)
		node.SetMeasureFunc(measureNode)
		node.StyleSetFlexShrink(1)
		node.StyleSetFlexGrow(0)
		node.StyleSetFlexDirection(flex.FlexDirectionRow)

		// If our node has layout properties, grab those.
		if el, ok := el.(ElementLayout); ok {
			el.Layout().apply(node)
		}

		// Setup our contxt
		elCache[idx] = &measureContext{
			Element: el,
		}
		node.Context = elCache[idx]

		// Insert our child
		root.InsertChild(node, idx)
	}

	// Calculate the layout
	flex.CalculateLayout(root, flex.Undefined, flex.Undefined, flex.DirectionLTR)

	if false {
		// Debug
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

	// Render each
	for idx, elCtx := range elCache {
		child := root.GetChild(idx)
		if child == nil {
			break
		}

		// If the height that the layout engine calculated is less than
		// the height that we originally measured, then we need to give the
		// element a chance to rerender into that height. If it still exceeds
		// it, we truncate.
		height := child.LayoutGetHeight()
		if height < elCtx.Size.Height {
			elCtx.Text = truncateTextHeight(elCtx.Text, int(height))
		}

		fmt.Fprint(d.w, elCtx.Text)
		if len(elCtx.Text) > 0 && elCtx.Text[len(elCtx.Text)-1] != '\n' {
			fmt.Fprintln(d.w)
		}
	}

	// Store how much we drew
	d.lastCount = uint(root.LayoutGetHeight())
}

var b = aec.EmptyBuilder
