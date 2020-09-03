package dynamiccli

// Elements are the individual items that are rendered within a document.
type Element interface {
	// Render is called to render this element.
	//
	// The rows/cols given are advisory. If the cols are ignored, the return
	// value may be wrapped or truncated (depending on layout settings). This
	// behavior may be undesirable and so it is recommended you remain within
	// the advisory amounts. If the rows are ignored, the output will be
	// truncated.
	Render(rows, cols uint) string
}

// ElementTerminalSizer can be implemented to receive the terminal size.
// See the function docs for more information.
type ElementTerminalSizer interface {
	Element

	// SetTerminalSize is called with the full terminal size. This may
	// exceed the size given by Render in certain cases. This will be called
	// before Render and Layout.
	SetTerminalSize(rows, cols uint)
}

type ElementLayout interface {
	Element

	Layout() *Layout
}
