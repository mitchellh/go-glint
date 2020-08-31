package dynamiccli

import (
	"io"
)

// Elements are the individual items that are rendered within a document.
type Element interface {
	// Render is called to render this element. This should NOT render a
	// trailing newline; the document itself will append a trailing newline
	// if necessary.
	//
	// The return value notes the number of lines that were drawn.  This
	// count includes the final line that doesn't end with a trailing newline.
	// It is very important that the number of lines are correct in any
	// implementation or rendering artifacts will occur.
	Render(w io.Writer, width uint) uint

	// Dynamic returns true if this element may change in the future.
	// If this is false then at some point the element rendering may be
	// complete and never called again. For example, text elements can be
	// non-dynamic (static): once the line is written, they don't need to be
	// written again.
	Dynamic() bool
}
