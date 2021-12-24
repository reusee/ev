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

var _ Sink = new(Dashboard)

func (d *Dashboard) Put(ev *Ev) error {
	d.scope.MutateCall(
		func(
			evs Evs,
		) *Evs {
			evs = append(evs, ev)
			return &evs
		},
	)
	d.scope.Call(func(
		refresh Refresh,
	) {
		refresh(d.screen)
	})
	return nil
}

func (d *Dashboard) Run() {
	for {

		d.scope.Call(func(
			refresh Refresh,
		) {
			refresh(d.screen)
		})

		ev := d.screen.PollEvent()
		switch ev := ev.(type) {

		case *tcell.EventResize:
			d.screen.Sync()

		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEsc || ev.Key() == tcell.KeyCtrlC {
				return
			}

		case *tcell.EventMouse:
			buttons := ev.Buttons()

			// scroll
			d.scope.MutateCall(func(
				evs Evs,
				offset EvsOffset,
			) any {

				if buttons&tcell.WheelUp > 0 {
					offset++
					if int(offset) > len(evs)-1 {
						offset = EvsOffset(len(evs) - 1)
					}
					return &offset
				}

				if buttons&tcell.WheelDown > 0 {
					_, height := d.screen.Size()
					offset--
					if offset < EvsOffset(-height+1) {
						offset = EvsOffset(-height + 1)
					}
					return &offset
				}

				return nil
			})

			// click
			d.scope.Call(func(
				handlers ClickHandlers,
			) {
				if buttons&tcell.Button1 == 0 {
					return
				}
				x, y := ev.Position()
				for _, handler := range handlers {
					if x >= handler.X && x <= handler.X+handler.Width &&
						y >= handler.Y && y <= handler.Y+handler.Height {
						if handler.Func != nil {
							d.scope.Call(handler.Func)
							break
						}
					}
				}
			})

		}

	}
}
