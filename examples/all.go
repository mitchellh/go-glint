package main

import (
	"context"
	"os"

	"github.com/mitchellh/go-dynamic-cli"
)

func main() {
	var d dynamiccli.Document
	d.SetOutput(os.Stdout)
	d.SetWidth(80)
	d.Add(
		dynamiccli.Text("Hello"),
		dynamiccli.Progress(100),
	)
	d.Render(context.Background())
}
