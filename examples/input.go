package main

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/mitchellh/go-glint"
)

func main() {
	var m sync.Mutex
	var buf bytes.Buffer

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d := glint.New()
	d.Append(
		glint.Text("Type things..."),
		glint.Input(func(data []byte, err error) {
			m.Lock()
			defer m.Unlock()
			buf.Write(data)

			for _, b := range data {
				if b == 3 {
					cancel()
				}
			}
		}),
		glint.TextFunc(func(height, width uint) string {
			m.Lock()
			defer m.Unlock()
			return fmt.Sprintf("%v", buf.Bytes())
		}),
	)
	d.Render(ctx)
}
