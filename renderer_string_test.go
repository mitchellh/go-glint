package glint

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringRenderer(t *testing.T) {
	require := require.New(t)

	r := &StringRenderer{}
	d := New()
	d.SetRenderer(r)
	d.Append(Text("hello\nworld"))

	d.RenderFrame()
	require.Equal("hello\nworld", r.Builder.String())

	// Second render should clear and rewrite
	d.RenderFrame()
	require.Equal("hello\nworld", r.Builder.String())
}
