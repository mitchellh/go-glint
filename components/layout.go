package components

import (
	"github.com/mitchellh/go-dynamic-cli"
	"github.com/mitchellh/go-dynamic-cli/internal/flex"
	"github.com/mitchellh/go-dynamic-cli/internal/layout"
)

func Layout(inner ...dynamiccli.Component) *LayoutComponent {
	return &LayoutComponent{inner: inner, builder: &layout.Builder{}}
}

type LayoutComponent struct {
	inner   []dynamiccli.Component
	builder *layout.Builder
}

func (c *LayoutComponent) Row() *LayoutComponent {
	c.builder = c.builder.Raw(func(n *flex.Node) {
		n.StyleSetFlexDirection(flex.FlexDirectionRow)
	})
	return c
}

func (c *LayoutComponent) MarginLeft(x int) *LayoutComponent {
	c.builder = c.builder.Raw(func(n *flex.Node) {
		n.StyleSetMargin(flex.EdgeLeft, float32(x))
	})
	return c
}

func (c *LayoutComponent) Body() dynamiccli.Component {
	return dynamiccli.Fragment(c.inner...)
}

func (c *LayoutComponent) Layout() *layout.Builder {
	return c.builder
}
