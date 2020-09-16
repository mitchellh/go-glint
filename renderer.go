package glint

import (
	"context"
	"io"

	"github.com/mitchellh/go-glint/flex"
)

// Renderers are responsible for helping configure layout properties and
// ultimately drawing components.
type Renderer interface {
	// LayoutRoot returns the root node for the layout engine. This should
	// set any styling to restrict children such as width. If this returns nil
	// then rendering will do nothing.
	LayoutRoot() *flex.Node

	// RenderRoot is called to render the tree rooted at the given node.
	// This will always be called with the root node. In the future we plan
	// to support partial re-renders but this will be done via a separate call.
	//
	// prev will be the previous root that was rendered. This can be used to
	// determine layout differences. This will be nil if this is the first
	// render.
	RenderRoot(root, prev *flex.Node)
}

// RendererInputStream can optionally be implemented by renderers that
// support an input stream. This is used by the Input component. If this
// isn't supported, Input will report that Inputs are not supported.
type RendererInputStream interface {
	Renderer

	// InputStream should return the io.Reader for input data. This may be
	// handled specially if this is a TTY. If you don't want a TTY to be
	// handled specially, wrap the reader in something like a bufio.Reader.
	InputStream() io.Reader
}

// WithRenderer inserts the renderer into the context. This is done automatically
// by Document for components.
func WithRenderer(ctx context.Context, r Renderer) context.Context {
	return context.WithValue(ctx, rendererCtxKey, r)
}

// RendererFromContext returns the Renderer in the context or nil if no
// Renderer is found.
func RendererFromContext(ctx context.Context) Renderer {
	v, _ := ctx.Value(rendererCtxKey).(Renderer)
	return v
}

type glintCtxKey string

const (
	rendererCtxKey = glintCtxKey("renderer")
)
