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

	exit := make(chan struct{})
	go func() {
		board.Run()
		close(exit)
	}()

	for i := 0; i < 128; i++ {
		put(&Ev{
			Name: "foo",
			Attrs: []Attr{
				{Name: "a1", Value: 42},
				{Name: "a2", Value: i},
			},
			Subs: Evs{
				{
					Name: "sub1",
					Attrs: []Attr{
						{Name: "1", Value: "yes"},
					},
				},
			},
		})
	}

	<-exit
}
