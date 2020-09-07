package main

import (
	"context"
	"strings"

	"github.com/mitchellh/go-glint"
)

func main() {
	d := glint.New()
	d.Append(
		glint.Text(strings.Repeat("123456", 12)),
		glint.Text(strings.Repeat("123456", 6)),
		glint.Text(strings.Repeat("123456", 3)),
	)
	d.Render(context.Background())
}
