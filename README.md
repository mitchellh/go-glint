# go-glint [![Godoc](https://godoc.org/github.com/mitchellh/go-glint?status.svg)](https://godoc.org/github.com/mitchellh/go-glint)

Glint is a component-based UI framework specifically targeted towards
command-line interfaces. This allows you to create highly dynamic CLI interfaces
using shared, easily testable components. Glint uses a Flexbox implementation
to make it easy to lay out components in the CLI, including paddings, margins,
and more.

**API Status: Unstable.** We're still actively working on the API and
may change it in backwards incompatible ways. See the roadmap section in
particular for work that may impact the API.

## Roadmap

Glint is still an early stage project and there is a lot that we want to
improve on. This may introduce some backwards incompatibilities but we are
trying to stabilize the API as quickly as possible.

* **Non-interactive interfaces.** We want to add support for rendering to
non-interactive interfaces and allowing components to provide custom behavior
in these cases. For now, users of Glint should detect non-interactivity and
avoid using Glint.

* **Windows PowerShell and Cmd.** Glint works fine in ANSI-compatible terminals
on Windows, but doesn't work with PowerShell and Cmd. We want to make this
work.

* **User Input.** Glint should be able to query for user input and render
this within its existing set of components.

* **Mount/Unmount Callbacks.** Glint components should be provided with a way
to detect when they're mounted, unmounted, and other events within a document.
This can be used for initialization and cleanup.

* **Custom renderers.** We want to support custom renderers that can take
the common component tree and choose to render and draw frames differently.
This would allow Glint to be able to output JSON streams, draw to interfaces
such as Slack, and others all using the same component model.
