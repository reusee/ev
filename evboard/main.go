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
	board, err := ev.NewDashboard(ctx)
	ce(err)
	defer board.Close()

	exit := make(chan struct{})
	go func() {
		board.Run()
		close(exit)
	}()

	errCh := make(chan error, 1)
	ce(ev.NewTCPServer(ctx, ":9876", board, errCh))

	for {
		select {
		case <-exit:
			return
		case err := <-errCh:
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				break
			}
			ce(err)
		}
	}
}
