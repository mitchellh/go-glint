package glint

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLayout(t *testing.T) {
	t.Run("left margin", func(t *testing.T) {
		// single line
		require.Equal(t, "  hello", TestRender(t,
			Layout(Text("hello")).MarginLeft(2),
		))
	})

	t.Run("right margin", func(t *testing.T) {
		// single line
		require.Equal(t, "hello  ", TestRender(t,
			Layout(Text("hello")).MarginRight(2),
		))
	})

	t.Run("left padding", func(t *testing.T) {
		// single line
		require.Equal(t, "  hello", TestRender(t,
			Layout(Text("hello")).PaddingLeft(2),
		))
	})

	t.Run("right padding", func(t *testing.T) {
		// single line
		require.Equal(t, "hello  ", TestRender(t,
			Layout(Text("hello")).PaddingRight(2),
		))
	})
}
