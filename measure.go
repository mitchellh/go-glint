package glint

import (
	"strings"

	"github.com/mitchellh/go-glint/internal/flex"
)

func measureNode(
	node *flex.Node,
	width float32,
	widthMode flex.MeasureMode,
	height float32,
	heightMode flex.MeasureMode,
) flex.Size {
	// If we have no context set then we use the full spacing.
	ctx, ok := node.Context.(*nodeContext)
	if !ok || ctx == nil {
		return flex.Size{Width: width, Height: height}
	}

	// Otherwise, we have to render this.
	ctx.Text = ctx.Component.render(uint(height), uint(width))
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
