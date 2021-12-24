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

	var sync DashboardSync
	d.scope.Assign(&sync)

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
	) Evs {
		return append(evs, ev)
	})
	var sync DashboardSync
	d.scope.Assign(&sync)
	return nil
}

type DashboardSync bool

func (_ DashboardDef) DashboardSync(
	evs Evs,
) DashboardSync {

	//TODO

	return true
}
