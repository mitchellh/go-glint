package glint

func Fragment(c ...Component) *fragmentComponent {
	return &fragmentComponent{List: c}
}

type fragmentComponent struct {
	terminalComponent

	List []Component
}
