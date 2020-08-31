package dynamiccli

import (
	"fmt"
	"io"
)

// TextElements is an Element that renders text.
type TextElement struct {
	text string
}

// Text creates a TextElement for static text. The text here will be word
// wrapped automatically based on the width of the terminal.
func Text(v string) *TextElement {
	return &TextElement{
		text: v,
	}
}

func (el *TextElement) Render(w io.Writer, width uint) uint {
	fmt.Fprint(w, el.text)
	return 1
}

func (el *TextElement) Dynamic() bool {
	return false
}
