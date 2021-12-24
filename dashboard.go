package ev

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/reusee/dscope"
)

type Dashboard struct {
	scope  *dscope.MutableScope
	screen tcell.Screen
}

func NewDashboard(ctx context.Context) (*Dashboard, error) {
	scope := dscope.NewMutable(dscope.Methods(new(DashboardDef))...)

	d := &Dashboard{
		scope: scope,
	}

	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, we(err)
	}
	if err := screen.Init(); err != nil {
		return nil, we(err)
	}
	screen.EnableMouse()
	screen.EnablePaste()
	screen.SetStyle(tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorReset))
	screen.Clear()
	d.screen = screen

	return d, nil
}

func (d *Dashboard) Close() {
	d.screen.Fini()
}

func (d *Dashboard) update(fn any) {
	d.scope.MutateCall(fn)
	d.refresh()
}

var _ Sink = new(Dashboard)

func (d *Dashboard) Put(ev *Ev) error {
	d.update(
		func(
			evs Evs,
		) *Evs {
			evs = append(evs, ev)
			return &evs
		},
	)
	return nil
}

func (d *Dashboard) refresh() {
	d.screen.Clear()
	d.scope.Call(func(
		refresh DashboardRefresh,
	) {
		refresh(d.screen)
	})
	d.screen.Show()
}

func (d *Dashboard) Run() {
	for {

		ev := d.screen.PollEvent()
		switch ev := ev.(type) {

		case *tcell.EventResize:
			d.update(func() {})
			d.screen.Sync()

		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEsc || ev.Key() == tcell.KeyCtrlC {
				return
			}

		case *tcell.EventMouse:
			d.update(func(
				evs Evs,
				offset EvsOffset,
			) *EvsOffset {
				buttons := ev.Buttons()
				if buttons&tcell.WheelUp > 0 {
					offset++
					if int(offset) > len(evs)-1 {
						offset = EvsOffset(len(evs) - 1)
					}
					return &offset
				} else if buttons&tcell.WheelDown > 0 {
					offset--
					if offset < 0 {
						offset = 0
					}
				}
				return &offset
			})

		}

	}
}
