package components

import (
	"github.com/mitchellh/go-glint"
	"github.com/tj/go-spin"
)

// Spinner creates a new spinner. The created spinner should NOT be started
// or data races will occur that can result in a panic.
func Spinner() *SpinnerComponent {
	// Create our spinner and setup our default frames
	s := spin.New()
	s.Set(spin.Default)

	return &SpinnerComponent{
		s: s,
	}
}

type SpinnerComponent struct {
	s *spin.Spinner
}

func (c *SpinnerComponent) Body() glint.Component {
	return glint.Text(c.s.Next())
}
