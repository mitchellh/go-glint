package glint

import (
	"sync"
)

// TextComponent is a Component that renders text.
type TextComponent struct {
	mu   sync.Mutex
	text string
	f    func(rows, cols uint) string
}

// Text creates a TextComponent for static text. The text here will be word
// wrapped automatically based on the width of the terminal.
func Text(v string) *TextComponent {
	return &TextComponent{
		text: v,
	}
}

func TextFunc(f func(rows, cols uint) string) *TextComponent {
	return &TextComponent{
		f: f,
	}
}

// Update updates the text element. This is safe to call while this is being
// rendered.
func (el *TextComponent) Update(text string) {
	el.mu.Lock()
	defer el.mu.Unlock()
	el.text = text
}

func (el *TextComponent) Body() Component {
	return nil
}

func (el *TextComponent) render(rows, cols uint) string {
	el.mu.Lock()
	defer el.mu.Unlock()

	if el.f != nil {
		return el.f(rows, cols)
	}

	return el.text
}
