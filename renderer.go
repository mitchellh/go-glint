package glint

import (
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
