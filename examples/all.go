package main

import (
	"context"
	"os"

	"github.com/mitchellh/go-dynamic-cli"
)

func main() {
	var d dynamiccli.Document
	d.SetOutput(os.Stdout)
	d.Add(
		dynamiccli.Text("Hello"),
	)
	d.Render(context.Background())
}
