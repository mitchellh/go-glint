package layout

import (
	"github.com/mitchellh/go-dynamic-cli/internal/flex"
)

type SetFunc func(n *flex.Node)

type Builder struct {
	f SetFunc
}

func (l *Builder) Raw(f SetFunc) *Builder {
	return l.add(f)
}

func (l *Builder) Apply(node *flex.Node) {
	if l == nil || l.f == nil {
		return
	}

	l.f(node)
}

func (l *Builder) add(f func(*flex.Node)) *Builder {
	old := l.f
	new := func(n *flex.Node) {
		if old != nil {
			old(n)
		}

		f(n)
	}

	return &Builder{f: new}
}
