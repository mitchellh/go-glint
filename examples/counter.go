package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/mitchellh/go-glint"
)

func main() {
	var counter uint32
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			atomic.AddUint32(&counter, 1)
		}
	}()

	d := glint.New()
	d.Append(
		glint.Style(
			glint.TextFunc(func(rows, cols uint) string {
				return fmt.Sprintf("%d tests passed", atomic.LoadUint32(&counter))
			}),
			glint.Color("green"),
		),
	)
	d.Render(context.Background())
}
