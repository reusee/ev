package ev

import (
	"context"
	"testing"

	"github.com/reusee/e4"
)

func TestDashboard(t *testing.T) {
	defer he(nil, e4.TestingFatal(t))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	board, err := NewDashboard(ctx)
	ce(err)
	defer board.Close()

	put := MustPut(board)
	_ = put
}
