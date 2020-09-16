package glint

import (
	"context"

	"github.com/mitchellh/go-glint/internal/input"
)

// InputFunc is the callback type for any input.
type InputFunc func([]byte, error)

// InputComponent is a Component that renders text.
type InputComponent struct {
	terminalComponent
	cb InputFunc
}

func Input(cb InputFunc) *InputComponent {
	return &InputComponent{cb: cb}
}

func (el *InputComponent) Mount(ctx context.Context) {
	// Get our renderer
	r, ok := RendererFromContext(ctx).(RendererInputStream)
	if !ok || r == nil {
		return
	}

	// Get our reader. If it is a TTY then we turn raw mode on.
	reader := r.InputStream()

	// Get our manager and add ourselves
	inputMgr := input.For(r, reader)
	inputMgr.AddCallback(el, input.Func(el.cb))
}

func (el *InputComponent) Unmount(ctx context.Context) {
	// Get our renderer
	r, ok := RendererFromContext(ctx).(RendererInputStream)
	if !ok || r == nil {
		return
	}

	// Get our manager and run a close
	input.Close(r)
}

func (el *InputComponent) Body(context.Context) Component {
	return nil
}
