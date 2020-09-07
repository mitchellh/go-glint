package glint

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/creack/pty"
	"github.com/morikuni/aec"
	sshterm "golang.org/x/crypto/ssh/terminal"

	"github.com/mitchellh/go-glint/internal/flex"
)

// TerminalRenderer renders output to a terminal. It expects the Output set
// to be a TTY. This will use ANSI escape codes to redraw.
type TerminalRenderer struct {
	// Output is where to write to. This should be a TTY.
	Output io.Writer

	// Rows, Cols are the dimensions of the terminal. If these are not set
	// (zero), then we will auto-detect the size of the output if it is a TTY.
	// If the values are still zero, nothing will be rendered.
	Rows, Cols uint
}

func (r *TerminalRenderer) LayoutRoot() *flex.Node {
	// If we don't have a writer set, then don't render anything.
	if r.Output == nil {
		return nil
	}

	// Setup our dimensions
	cols := r.Cols
	rows := r.Rows
	if cols == 0 || rows == 0 {
		if f, ok := r.Output.(*os.File); ok && sshterm.IsTerminal(int(f.Fd())) {
			ws, err := pty.GetsizeFull(f)
			if err == nil {
				rows = uint(ws.Rows)
				cols = uint(ws.Cols)
			}
		}
	}

	// Render nothing if we're going to have any zero dimensions
	if cols == 0 || rows == 0 {
		return nil
	}

	// Setup our node
	node := flex.NewNode()
	node.StyleSetWidth(float32(cols))
	node.Context = &termRootContext{
		Rows: rows,
		Cols: cols,
	}

	return node
}

func (r *TerminalRenderer) RenderRoot(root, prev *flex.Node) {
	w := r.Output
	rootCtx := root.Context.(*termRootContext)
	rows := rootCtx.Rows

	// Remove what we last drew. If what we last drew is greater than the number
	// of rows then we need to clear the screen.
	if prev != nil {
		height := uint(root.LayoutGetHeight())
		if height > 0 {
			if height <= rows {
				// Delete current line
				fmt.Fprint(w, b.Column(0).EraseLine(aec.EraseModes.All).ANSI)

				// Delete n lines above
				for i := uint(0); i < height-1; i++ {
					fmt.Fprint(w, b.Up(1).Column(0).EraseLine(aec.EraseModes.All).ANSI)
				}
			} else {
				fmt.Fprint(w, b.EraseDisplay(aec.EraseModes.All).EraseDisplay(aec.EraseMode(3)).Position(0, 0).ANSI)
			}
		}
	}

	// Draw
	r.renderTree(root, -1)
}

func (r *TerminalRenderer) renderTree(parent *flex.Node, lastRow int) {
	for _, child := range parent.Children {
		// If we're on a different row than last time then we draw a newline.
		thisRow := int(child.LayoutGetTop())
		if lastRow >= 0 && thisRow > lastRow {
			fmt.Fprintln(r.Output)
		}
		lastRow = thisRow

		// If we have a left margin, draw that first.
		if v := int(child.LayoutGetMargin(flex.EdgeLeft)); v > 0 {
			fmt.Fprint(r.Output, strings.Repeat(" ", v))
		}
		if v := int(child.LayoutGetPadding(flex.EdgeLeft)); v > 0 {
			fmt.Fprint(r.Output, strings.Repeat(" ", v))
		}

		// Get our node context. If we don't have one then we're a container
		// and we render below.
		ctx, ok := child.Context.(*TextNodeContext)
		if !ok {
			r.renderTree(child, lastRow)
		} else {
			// Draw our text
			fmt.Fprint(r.Output, ctx.Text)
		}

		// If we have a left margin, draw that first.
		if v := int(child.LayoutGetMargin(flex.EdgeRight)); v > 0 {
			fmt.Fprint(r.Output, strings.Repeat(" ", v))
		}
		if v := int(child.LayoutGetPadding(flex.EdgeRight)); v > 0 {
			fmt.Fprint(r.Output, strings.Repeat(" ", v))
		}
	}
}

type termRootContext struct {
	Rows, Cols uint
}

type termLeafContext struct {
	Component *TextComponent
	Text      string
	Size      flex.Size
}

var b = aec.EmptyBuilder
