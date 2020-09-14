package glint

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDocument_mountUnmount(t *testing.T) {
	require := require.New(t)

	// Create our doc
	r := &StringRenderer{}
	d := New()
	d.SetRenderer(r)

	// Add our component
	var c testMount
	d.Append(&c)
	require.Equal(uint32(0), atomic.LoadUint32(&c.mount))
	require.Equal(uint32(0), atomic.LoadUint32(&c.unmount))

	// Render once
	d.RenderFrame()
	require.Equal(uint32(1), atomic.LoadUint32(&c.mount))
	require.Equal(uint32(0), atomic.LoadUint32(&c.unmount))

	// Render again
	d.RenderFrame()
	require.Equal(uint32(1), atomic.LoadUint32(&c.mount))
	require.Equal(uint32(0), atomic.LoadUint32(&c.unmount))

	// Remove the old components
	d.Set()
	d.RenderFrame()
	require.Equal(uint32(1), atomic.LoadUint32(&c.mount))
	require.Equal(uint32(1), atomic.LoadUint32(&c.unmount))

	// Render again
	d.RenderFrame()
	require.Equal(uint32(1), atomic.LoadUint32(&c.mount))
	require.Equal(uint32(1), atomic.LoadUint32(&c.unmount))
}

type testMount struct {
	terminalComponent

	mount   uint32
	unmount uint32
}

func (c *testMount) Mount()   { atomic.AddUint32(&c.mount, 1) }
func (c *testMount) Unmount() { atomic.AddUint32(&c.unmount, 1) }

var (
	_ Component        = (*testMount)(nil)
	_ ComponentMounter = (*testMount)(nil)
)
