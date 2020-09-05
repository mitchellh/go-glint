package main

import (
	"context"
	"os"
	"time"

	"github.com/mitchellh/go-glint"
	gc "github.com/mitchellh/go-glint/components"
)

func main() {
	var d glint.Document
	d.SetOutput(os.Stdout)
	d.Append(
		gc.Layout(
			gc.Spinner(),
			gc.Layout(glint.Text("Build site and validate links...")).MarginLeft(1),
			gc.Layout(gc.Stopwatch(time.Now())).MarginLeft(1),
		).Row(),
		gc.Layout(
			gc.Spinner(),
			gc.Layout(glint.Text("Preparing execution environment...")).MarginLeft(1),
			gc.Layout(gc.Stopwatch(time.Now())).MarginLeft(1),
		).MarginLeft(2).Row(),
		gc.Layout(
			glint.Text("Preparing volume to work with..."),
		).MarginLeft(4),
	)
	d.Render(context.Background())
}
