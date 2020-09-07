package glint

import (
	"gopkg.in/gookit/color.v1"
)

// TextComponent is a Component that renders text.
type TextComponent struct {
	terminalComponent
	f func(rows, cols uint) string

	fgColor, bgColor colorizer
	style            []color.Color
}

// Text creates a TextComponent for static text. The text here will be word
// wrapped automatically based on the width of the terminal.
func Text(v string, opts ...TextOption) *TextComponent {
	return TextFunc(func(rows, cols uint) string { return v }, opts...)
}

// TextFunc creates a TextComponent for text that is dependent on the
// size of the draw area.
func TextFunc(f func(rows, cols uint) string, opts ...TextOption) *TextComponent {
	c := &TextComponent{
		f: f,
	}
	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (el *TextComponent) Body() Component {
	return nil
}

func (el *TextComponent) Render(rows, cols uint) string {
	if el.f == nil {
		return ""
	}

	return el.f(rows, cols)
}

// colorize colors the string using terminal escape codes according to the
// color set on text. This is purposely NOT exported because I want to
// expose colors to custom renderers in the future in a more consumable way.
func (el *TextComponent) colorize(v string) string {
	// Set colors
	if el.bgColor != nil {
		v = el.bgColor.Sprint(v)
	}
	if el.fgColor != nil {
		v = el.fgColor.Sprint(v)
	}
	v = color.Style(el.style).Sprint(v)

	return v
}

// TextOption is an option that can be set when creating Text components.
type TextOption func(t *TextComponent)

// Color sets the color by name. The supported colors are listed below.
//
// black, red, green, yellow, blue, magenta, cyan, white, darkGray,
// lightRed, lightGreen, lightYellow, lightBlue, lightMagenta, lightCyan,
// lightWhite.
func Color(name string) TextOption {
	return func(t *TextComponent) {
		if c, ok := color.FgColors[name]; ok {
			t.fgColor = c
		}
		if c, ok := color.ExFgColors[name]; ok {
			t.fgColor = c
		}
	}
}

// ColorHex sets the foreground color by hex code. The value can be
// in formats AABBCC, #AABBCC, 0xAABBCC.
func ColorHex(v string) TextOption {
	return func(t *TextComponent) {
		t.fgColor = color.HEX(v)
	}
}

// ColorRGB sets the foreground color by RGB values.
func ColorRGB(r, g, b uint8) TextOption {
	return func(t *TextComponent) {
		t.fgColor = color.RGB(r, g, b)
	}
}

// BGColor sets the color by name. The supported colors are listed below.
//
// black, red, green, yellow, blue, magenta, cyan, white, darkGray,
// lightRed, lightGreen, lightYellow, lightBlue, lightMagenta, lightCyan,
// lightWhite.
func BGColor(name string) TextOption {
	return func(t *TextComponent) {
		if c, ok := color.BgColors[name]; ok {
			t.bgColor = c
		}
		if c, ok := color.ExBgColors[name]; ok {
			t.bgColor = c
		}
	}
}

// BGColorHex sets the foreground color by hex code. The value can be
// in formats AABBCC, #AABBCC, 0xAABBCC.
func BGColorHex(v string) TextOption {
	return func(t *TextComponent) {
		t.bgColor = color.HEX(v, true)
	}
}

// BGColorRGB sets the foreground color by RGB values.
func BGColorRGB(r, g, b uint8) TextOption {
	return func(t *TextComponent) {
		t.bgColor = color.RGB(r, g, b, true)
	}
}

// Bold sets the text to bold.
func Bold() TextOption {
	return func(t *TextComponent) {
		t.style = append(t.style, color.OpBold)
	}
}

// Italic sets the text to italic.
func Italic() TextOption {
	return func(t *TextComponent) {
		t.style = append(t.style, color.OpItalic)
	}
}

// Underline sets the text to be underlined.
func Underline() TextOption {
	return func(t *TextComponent) {
		t.style = append(t.style, color.OpUnderscore)
	}
}

type colorizer interface {
	Sprint(...interface{}) string
}
