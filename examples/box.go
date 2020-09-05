package main

import (
	"context"
	"os"
	"time"

	"github.com/mitchellh/go-glint"
	dc "github.com/mitchellh/go-glint/components"
)

func main() {
	var d glint.Document
	d.SetOutput(os.Stdout)
	d.Append(
		dc.Layout(
			dc.Spinner(),
			dc.Layout(glint.Text("Build site and validate links...")).MarginLeft(1),
			dc.Layout(dc.Stopwatch(time.Now())).MarginLeft(1),
		).Row(),
		dc.Layout(
			dc.Spinner(),
			dc.Layout(glint.Text("Preparing execution environment...")).MarginLeft(1),
			dc.Layout(dc.Stopwatch(time.Now())).MarginLeft(1),
		).MarginLeft(2).Row(),
		dc.Layout(
			glint.Text("Preparing volume to work with..."),
		).MarginLeft(4),
	)
	d.Render(context.Background())
}
