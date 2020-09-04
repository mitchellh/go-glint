package dynamiccli

// Finalize reutrns a component that will finalize the input component.
func Finalize(c Component) Component {
	return &finalizedComponent{
		Component: c,
	}
}

type finalizedComponent struct {
	Component
}

func (c *finalizedComponent) Body() Component {
	return c.Component
}
