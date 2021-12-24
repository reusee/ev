package ev

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	lru "github.com/hashicorp/golang-lru"
)

type DashboardDef struct{}

func (_ DashboardDef) Evs() Evs {
	return nil
}

type EvsOffset int

func (_ DashboardDef) EvsOffset() EvsOffset {
	return 0
}

type DashboardRefresh func(screen tcell.Screen)

func (_ DashboardDef) DashboardRefresh(
	evs Evs,
	offset EvsOffset,
	evBox EvBox,
) DashboardRefresh {
	return func(screen tcell.Screen) {

		width, height := screen.Size()
		x := 0
		y := height - 1

		// evs
		for i := len(evs) - 1 - int(offset); i >= 0; i-- {
			box := evBox(evs[i])
			y += box.DrawLeftBottom(screen, x, y, width, height)
			if y < 0 {
				break
			}
		}

	}
}

type Styles struct {
	EvName   tcell.Style
	AttrName tcell.Style
}

func (_ DashboardDef) Styles() (styles Styles) {
	styles.EvName = tcell.StyleDefault.Foreground(tcell.ColorRed)
	styles.AttrName = tcell.StyleDefault.Foreground(tcell.ColorGreen)
	return
}

type EvBox func(ev *Ev) BoxSpec

func (_ DashboardDef) EvBox(
	styles Styles,
) EvBox {
	cache, err := lru.New(512)
	ce(err)

	var evBox func(ev *Ev, level int) []RowSpec
	evBox = func(ev *Ev, level int) (rows []RowSpec) {

		var row RowSpec
		row.AppendStr(strings.Repeat("  ", level), tcell.StyleDefault)
		row.AppendStr(ev.Name, styles.EvName)
		for _, attr := range ev.Attrs {
			row.AppendStr(" ", tcell.StyleDefault)
			row.AppendStr(attr.Name, styles.AttrName)
			row.AppendStr(": ", tcell.StyleDefault)
			row.AppendStr(fmt.Sprintf("%v", attr.Value), tcell.StyleDefault)
		}
		rows = append(rows, row)

		for _, sub := range ev.Subs {
			rows = append(rows,
				evBox(sub, level+1)...)
		}

		return
	}

	return func(ev *Ev) BoxSpec {
		if v, ok := cache.Get(ev); ok {
			return v.(BoxSpec)
		}

		rows := evBox(ev, 0)
		box := BoxSpec{
			Rows: rows,
		}
		cache.Add(ev, box)

		return box
	}
}
