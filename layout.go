package dynamiccli

import (
	"strings"

	"github.com/mitchellh/go-dynamic-cli/internal/flex"
)

type Layout struct {
	f func(*flex.Node)
}

func NewLayout() *Layout {
	return &Layout{}
}

func (l *Layout) apply(node *flex.Node) {
	if l == nil || l.f == nil {
		return
	}

	l.f(node)
}

func (l *Layout) add(f func(*flex.Node)) *Layout {
	old := l.f
	new := func(n *flex.Node) {
		if old != nil {
			old(n)
		}

		f(n)
	}

	return &Layout{f: new}
}

func (l *Layout) MinHeight(v float32) *Layout {
	return l.add(func(n *flex.Node) {
		n.StyleSetMinHeight(v)
	})
}

func (l *Layout) Overflow(v Overflow) *Layout {
	return l.add(func(n *flex.Node) {
		n.StyleSetOverflow(flex.Overflow(v))
	})
}

type Overflow int

const (
	// OverflowVisible is "visible"
	OverflowVisible Overflow = iota
	// OverflowHidden is "hidden"
	OverflowHidden
	// OverflowScroll is "scroll"
	OverflowScroll
)

type measureContext struct {
	Element Element
	Text    string
	Size    flex.Size
}

func measureNode(
	node *flex.Node,
	width float32,
	widthMode flex.MeasureMode,
	height float32,
	heightMode flex.MeasureMode,
) flex.Size {
	// If we have no context set then we use the full spacing.
	ctx, ok := node.Context.(*measureContext)
	if !ok || ctx == nil {
		return flex.Size{Width: width, Height: height}
	}

	// Otherwise, we have to render this.
	ctx.Text = ctx.Element.Render(uint(height), uint(width))
	ctx.Size = flex.Size{
		Width:  float32(longestLine(ctx.Text)),
		Height: float32(countLines(ctx.Text)),
	}

	// TODO: wrapping

	return ctx.Size
}

func countLines(s string) int {
	count := strings.Count(s, "\n")

	// If the last character isn't a newline, we have to add one since we'll
	// always have one more line than newline characters.
	if len(s) > 0 && s[len(s)-1] != '\n' {
		count++
	}
	return count
}

func longestLine(s string) int {
	lastIdx, longest := 0, 0
	for {
		idx := strings.IndexByte(s, '\n')
		if idx == -1 {
			break
		}

		current := idx - lastIdx
		if current > longest {
			longest = current
		}

		s = s[idx+1:]
	}

	if longest == 0 {
		return len(s)
	}

	return longest
}

func truncateTextHeight(s string, height int) string {
	// The way this works is that we iterate through HEIGHT newlines
	// and return up to that point. If we either don't find a newline
	// or we've reached the end of the string, then the string is shorter
	// than the height limit and we return the whole thing.
	idx := 0
	for i := 0; i < height; i++ {
		next := strings.IndexByte(s[idx:], '\n')
		if next == -1 || idx >= len(s) {
			return s
		}

		idx += next + 1
	}

	if idx == 0 {
		return ""
	}

	// Subtract one here because the idx is the last "\n" char
	return s[:idx-1]
}
