package main

import (
	"context"
	"time"

	"github.com/mitchellh/go-glint"
	gc "github.com/mitchellh/go-glint/components"
)

func main() {
	d := glint.New()
	d.Append(
		glint.Layout(
			gc.Spinner(),
			glint.Layout(glint.Text("Build site and validate links...")).MarginLeft(1),
			glint.Layout(gc.Stopwatch(time.Now())).MarginLeft(1),
		).Row(),
		glint.Layout(
			gc.Spinner(),
			glint.Layout(glint.Text("Preparing execution environment...")).MarginLeft(1),
			glint.Layout(gc.Stopwatch(time.Now())).MarginLeft(1),
		).MarginLeft(2).Row(),
		glint.Layout(
			glint.Text("Preparing volume to work with..."),
		).MarginLeft(4),
	)
	d.Render(context.Background())
}
