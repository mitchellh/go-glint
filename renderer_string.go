package glint

import (
	"fmt"
	"strings"

	"github.com/mitchellh/go-glint/internal/flex"
)

// StringRenderer renders output to a string builder. This will clear
// the builder on each frame render. The StringRenderer is primarily meant
// for testing components.
type StringRenderer struct {
	// Builder is the strings builder to write to. If this is nil then
	// it will be created on first render.
	Builder *strings.Builder

	// Width is a fixed width to set for the root node. If this isn't
	// set then a width of 80 is arbitrarily used.
	Width uint
}

func (r *StringRenderer) LayoutRoot() *flex.Node {
	width := r.Width
	if width == 0 {
		width = 80
	}

	node := flex.NewNode()
	node.StyleSetWidth(float32(width))
	return node
}

func (r *StringRenderer) RenderRoot(root, prev *flex.Node) {
	if r.Builder == nil {
		r.Builder = &strings.Builder{}
	}

	// Reset our builder
	r.Builder.Reset()

	// Draw
	r.renderTree(root, -1)
}

func (r *StringRenderer) renderTree(parent *flex.Node, lastRow int) {
	for _, child := range parent.Children {
		// If we're on a different row than last time then we draw a newline.
		thisRow := int(child.LayoutGetTop())
		if lastRow >= 0 && thisRow > lastRow {
			r.Builder.WriteByte('\n')
		}
		lastRow = thisRow

		// If we have a left margin, draw that first.
		if v := int(child.LayoutGetMargin(flex.EdgeLeft)); v > 0 {
			fmt.Fprint(r.Builder, strings.Repeat(" ", v))
		}

		// Get our node context. If we don't have one then we're a container
		// and we render below.
		ctx, ok := child.Context.(*TextNodeContext)
		if !ok {
			r.renderTree(child, lastRow)
			continue
		}

		// Draw our text
		fmt.Fprint(r.Builder, ctx.Text)
	}
}
