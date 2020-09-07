package glint

import (
	"strings"

	"github.com/mitchellh/go-glint/internal/flex"
)

// TextNodeContext is the *flex.Node.Context set for all *TextComponent flex nodes.
type TextNodeContext struct {
	// C is the TextComponent represented.
	C *TextComponent

	// Text is the rendered text. This is populated after MeasureTextNode
	// is called. Note that this may not fit in the final layout calculations
	// since it is populated on measurement.
	Text string

	// Size is the measurement size returned. This can be used to determine
	// if the text above fits in the final size. Text is guaranteed to fit
	// in this size.
	Size flex.Size
}

// MeasureTextNode implements flex.MeasureFunc and returns the measurements
// for the given node only if the node represents a TextComponent. This is
// the MeasureFunc that is typically used for renderers since all component
// trees terminate in a text node.
//
// The flex.Node must have Context set to TextNodeContext. After calling this,
// fields such as Text and Size will be populated on the node.
func MeasureTextNode(
	node *flex.Node,
	width float32,
	widthMode flex.MeasureMode,
	height float32,
	heightMode flex.MeasureMode,
) flex.Size {
	// If we have no context set then we use the full spacing.
	ctx, ok := node.Context.(*TextNodeContext)
	if !ok || ctx == nil {
		return flex.Size{Width: width, Height: height}
	}

	// Otherwise, we have to render this.
	ctx.Text = ctx.C.Render(uint(height), uint(width))
	ctx.Size = flex.Size{
		Width:  float32(longestLine(ctx.Text)),
		Height: float32(countLines(ctx.Text)),
	}

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

	// This can happen if height == 0
	if idx == 0 {
		return ""
	}

	// Subtract one here because the idx is the last "\n" char
	return s[:idx-1]
}
