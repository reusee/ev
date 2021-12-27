package main

import (
	"context"
	"errors"
	"io"

	"github.com/reusee/ev"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	ce(ev.NewTCPServer(ctx, ":9876", ev.PutFunc(func(ev *ev.Ev) error {
		pt("%#v\n", ev)
		return nil
	}), errCh))

	for {
		select {
		case err := <-errCh:
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				break
			}
			ce(err)
		}
	}
}
