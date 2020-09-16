package glint

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/mitchellh/go-glint/flex"
)

// Document is the primary structure for managing and drawing components.
//
// A document represents a terminal window or session. The output can be set and
// components can be added, rendered, and drawn. All the methods on a Document
// are thread-safe unless otherwise documented. This allows you to draw,
// add components, replace components, etc. all while the render loop is active.
//
// Currently, this can only render directly to an io.Writer that expects to
// be a terminal session. In the future, we'll further abstract the concept
// of a "renderer" so that rendering can be done to other mediums as well.
type Document struct {
	mu          sync.Mutex
	r           Renderer
	els         []Component
	refreshRate time.Duration
	prevRoot    *flex.Node
	mounted     map[ComponentMounter]struct{}
}

// New returns a Document that will output to stdout.
func New() *Document {
	var d Document
	d.SetRenderer(&TerminalRenderer{
		Output: os.Stdout,
		Input:  os.Stdin,
	})

	return &d
}

// SetRenderer sets the renderer to use. If this isn't set then Render
// will do nothing and return immediately. Changes to this will have no
// impact on active render loops.
func (d *Document) SetRenderer(r Renderer) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.r = r
}

// SetRefreshRate sets the rate at which output is rendered.
func (d *Document) SetRefreshRate(dur time.Duration) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.refreshRate = dur
}

// Append appends components to the document.
func (d *Document) Append(el ...Component) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.els = append(d.els, el...)
}

// Set sets the components for the document. This will replace all
// previous components.
func (d *Document) Set(els ...Component) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.els = els
}

// Close ensures that all elements are unmounted by finalizing all the
// output and then calling RenderFrame. Users of Document should ensure
// that Close is always called.
func (d *Document) Close() error {
	d.mu.Lock()
	for i, el := range d.els {
		d.els[i] = Finalize(el)
	}
	d.mu.Unlock()

	// We call RenderFrame twice to ensure we remove the elements AND
	// call Unmount on them.
	d.RenderFrame()
	d.RenderFrame()

	return nil
}

// Render starts a render loop that continues to render until the
// context is cancelled. This will render at the configured refresh rate.
// If the refresh rate is changed, it will not affect an active render loop.
// You must cancel and restart the render loop.
func (d *Document) Render(ctx context.Context) {
	d.mu.Lock()
	dur := d.refreshRate
	d.mu.Unlock()
	if dur == 0 {
		dur = time.Second / 24
	}

	t := time.NewTicker(dur)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-t.C:
			d.RenderFrame()
		}
	}
}

// RenderFrame will render a single frame and return.
//
// If a manual size is not configured, this will recalcualte the window
// size on each call. This typically requires a syscall. This is a bit
// expensive but something we can optimize in the future if it ends up being
// a real source of FPS issues.
func (d *Document) RenderFrame() {
	d.mu.Lock()
	defer d.mu.Unlock()

	// If we don't have a renderer set, then don't render anything.
	if d.r == nil {
		return
	}

	// Our context
	ctx := WithRenderer(context.Background(), d.r)

	// Setup our root node
	root := d.r.LayoutRoot()

	// Build our render tree
	tree(ctx, root, Fragment(d.els...), false)

	// Calculate the layout
	flex.CalculateLayout(root, flex.Undefined, flex.Undefined, flex.DirectionLTR)

	// Fix any text nodes that need to be fixed.
	d.handleNodes(ctx, root, nil)

	// Render the tree
	d.r.RenderRoot(root, d.prevRoot)

	// Store how much we drew
	height := uint(root.LayoutGetHeight())

	// If our component list is prefixed with finalized components, we
	// prune these out and do not re-render them.
	finalIdx := -1
	for i, el := range d.els {
		child := root.GetChild(i)
		if child == nil {
			break
		}

		// If the component is not finalized then we exit. If the
		// component doesn't match our expectations it means we hit
		// something weird and we exit too.
		ctx, ok := child.Context.(*parentContext)
		if !ok || ctx == nil || ctx.C != el || !ctx.Finalized {
			break
		}

		// If this is finalized, then we have to subtract from the
		// height the height of this child since we're not going to redraw.
		// Then continue until we find one that isn't finalized.
		height -= uint(child.LayoutGetHeight())
		finalIdx = i
	}
	if finalIdx >= 0 {
		// Change our elements
		els := d.els[finalIdx+1:]
		d.els = make([]Component, len(els))
		copy(d.els, els)

		// Reset the height on the root so that it reflects this change
		root.Layout.Dimensions[flex.DimensionHeight] = float32(height)
	}

	// Store our previous root
	d.prevRoot = root
}

func (d *Document) handleNodes(
	ctx context.Context,
	parent *flex.Node,
	seen map[ComponentMounter]struct{},
) {
	// For our first call, we detect the root since we use it later
	// to do some final calls.
	root := seen == nil
	if root {
		seen = map[ComponentMounter]struct{}{}
	}

	for _, child := range parent.Children {
		if tctx, ok := child.Context.(treeContext); ok {
			c := tctx.Component()

			// Mount callbacks
			if mc, ok := c.(ComponentMounter); ok {
				// Only if we haven't seen this already...
				if _, ok := seen[mc]; !ok {
					seen[mc] = struct{}{}

					if d.mounted == nil {
						d.mounted = map[ComponentMounter]struct{}{}
					}

					// And we haven't notified this already...
					if _, ok := d.mounted[mc]; !ok {
						d.mounted[mc] = struct{}{}

						// Notify
						mc.Mount(ctx)
					}
				}
			}

			continue
		}

		// If the height/width that the layout engine calculated is less than
		// the height that we originally measured, then we need to give the
		// element a chance to rerender into that dimension.
		if tctx, ok := child.Context.(*TextNodeContext); ok {
			height := child.LayoutGetHeight()
			width := child.LayoutGetWidth()
			if height < tctx.Size.Height || width < tctx.Size.Width {
				child.Measure(child,
					width, flex.MeasureModeAtMost,
					height, flex.MeasureModeAtMost,
				)
			}
		}

		d.handleNodes(ctx, child, seen)
	}

	// If we're the root call, then we preform some final calls. Otherwise
	// we just return, we're done.
	if !root {
		return
	}

	// Go through our previously mounted set and if we didn't see it,
	// then call unmount on it. After we're done, what we saw is our new
	// map of mounted elements.
	for mc := range d.mounted {
		if _, ok := seen[mc]; !ok {
			mc.Unmount(ctx)
		}
	}
	d.mounted = seen
}
