package dynamiccli

import (
	"strings"

	"github.com/cheggaaa/pb/v3"
)

// ProgressElement renders a progress bar. This wraps the cheggaaa/pb package
// since that provides important functionality. This uses single call renders
// to render the progress bar as values change.
type ProgressElement struct {
	*pb.ProgressBar
}

// Progress creates a new progress bar element with the given total.
// For more fine-grained control, please construct a ProgressElement
// directly.
func Progress(total int) *ProgressElement {
	return &ProgressElement{
		ProgressBar: pb.New(total),
	}
}

func (el *ProgressElement) Render(rows, cols uint) string {
	// If we have no progress bar render nothing.
	if el.ProgressBar == nil {
		return ""
	}

	// Set the width so we render properly
	el.ProgressBar.SetWidth(int(cols))

	// Write the current progress
	//return el.ProgressBar.String()
	return strings.TrimSpace(strings.Repeat(el.ProgressBar.String()+"\n", 15))
}

func (el *ProgressElement) Layout() *Layout {
	return NewLayout().MinHeight(0).Overflow(OverflowHidden)
}
