package dynamiccli

import (
	"fmt"
	"io"

	"github.com/mitchellh/go-dynamic-cli/internal/flex"
)

func tree(
	parent *flex.Node,
	c Component,
	termRows, termCols uint,
) {
	// Don't do anything with no component
	if c == nil {
		return
	}

	// Setup our node
	node := flex.NewNodeWithConfig(parent.Config)
	parent.InsertChild(node, len(parent.Children))

	// Notify of the terminal size
	if c, ok := c.(ComponentTerminalSizer); ok {
		c.SetTerminalSize(termRows, termCols)
	}

	// Setup a custom layout
	if c, ok := c.(ComponentLayout); ok {
		c.Layout().apply(node)
	}

	switch c := c.(type) {
	case *fragmentComponent:
		for _, c := range c.List {
			tree(parent, c, termRows, termCols)
		}

	case *TextComponent:
		// If this is a terminal node then we setup extra styles
		node.Context = &nodeContext{
			Component: c,
		}

		node.StyleSetFlexShrink(1)
		node.StyleSetFlexGrow(0)
		node.StyleSetFlexDirection(flex.FlexDirectionRow)
		node.SetMeasureFunc(measureNode)

	default:
		// If this is not terminal then we nest.
		tree(node, c.Body(), termRows, termCols)
	}

}

func renderTree(w io.Writer, parent *flex.Node) {
	for _, child := range parent.Children {
		ctx, ok := child.Context.(*nodeContext)
		if !ok {
			renderTree(w, child)
			continue
		}

		text := ctx.Text

		// If the height/width that the layout engine calculated is less than
		// the height that we originally measured, then we need to give the
		// element a chance to rerender into that dimension. If it still exceeds
		// it, we truncate.
		height := child.LayoutGetHeight()
		width := child.LayoutGetWidth()
		if height < ctx.Size.Height || width < ctx.Size.Width {
			// Rerender into it
			text = ctx.Component.render(uint(height), uint(width))

			// Truncate, no-ops if it fits.
			text = truncateTextHeight(text, int(height))
		}

		fmt.Fprint(w, text)

		// If the text didn't end with a newline then we add one since
		// all elements here are block-level.
		if len(text) > 0 && text[len(text)-1] != '\n' {
			fmt.Fprintln(w)
		}
	}
}

type nodeContext struct {
	Component *TextComponent
	Text      string
	Size      flex.Size
}
