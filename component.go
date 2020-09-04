package dynamiccli

import (
	"github.com/mitchellh/go-dynamic-cli/internal/layout"
)

// Components are the individual items that are rendered within a document.
type Component interface {
	// Body returns the body of this component. This can be another custom
	// component or a standard component such as Text.
	Body() Component
}

// ComponentTerminalSizer can be implemented to receive the terminal size.
// See the function docs for more information.
type ComponentTerminalSizer interface {
	Component

	// SetTerminalSize is called with the full terminal size. This may
	// exceed the size given by Render in certain cases. This will be called
	// before Render and Layout.
	SetTerminalSize(rows, cols uint)
}

// componentLayout can be implemented to set custom layout settings
// for the component. This can only be implemented by internal components
// since we use an internal library.
//
// End users should use the "Layout" component to set layout options.
type componentLayout interface {
	Component

	// Layout should return the layout settings for this component.
	Layout() *layout.Builder
}

// terminalComponent is an embeddable struct for internal usage that
// satisfies Component. This is used since terminal components are handled
// as special cases.
type terminalComponent struct{}

func (terminalComponent) Body() Component { return nil }
