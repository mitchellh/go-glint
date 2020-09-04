package main

import (
	"context"
	"os"
	"time"

	"github.com/mitchellh/go-dynamic-cli"
	dc "github.com/mitchellh/go-dynamic-cli/components"
)

func main() {
	var d dynamiccli.Document
	d.SetOutput(os.Stdout)
	d.Add(
		dynamiccli.Finalize(dc.Stopwatch(time.Now())),
		dc.Layout(
			dc.Spinner(),
			dc.Layout(dynamiccli.Text("Build site and validate links...")).MarginLeft(1),
			dc.Layout(dc.Stopwatch(time.Now())).MarginLeft(1),
		).Row(),
		dc.Layout(
			dc.Spinner(),
			dc.Layout(dynamiccli.Text("Preparing execution environment...")).MarginLeft(1),
			dc.Layout(dc.Stopwatch(time.Now())).MarginLeft(1),
		).MarginLeft(2).Row(),
		dc.Layout(
			dynamiccli.Text("Preparing volume to work with..."),
		).MarginLeft(4),
	)
	d.Render(context.Background())
}
