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
	screen.SetStyle(tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset))
	screen.Clear()
	d.screen = screen

	return d, nil
}

func (d *Dashboard) Close() {
	d.screen.Fini()
}

type DashboardDef struct{}

func (_ DashboardDef) Evs() Evs {
	return nil
}

var _ Sink = new(Dashboard)

func (d *Dashboard) Put(ev *Ev) error {
	d.scope.MutateCall(func(
		evs Evs,
	) *Evs {
		evs = append(evs, ev)
		return &evs
	})
	d.refresh()
	return nil
}

func (d *Dashboard) refresh() {
	d.screen.Show()
}

func (d *Dashboard) Run() {
	for {

		ev := d.screen.PollEvent()
		switch ev := ev.(type) {

		case *tcell.EventResize:
			d.screen.Sync()

		case *tcell.EventKey:

			if ev.Key() == tcell.KeyEsc || ev.Key() == tcell.KeyCtrlC {
				return
			}

		}

	}
}
