package ev

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	lru "github.com/hashicorp/golang-lru"
	"github.com/reusee/dscope"
)

type DashboardDef struct{}

type Refresh func(screen tcell.Screen)

func (_ DashboardDef) Refresh(
	refresh DashboardRefresh,
) Refresh {
	var once sync.Once
	return func(screen tcell.Screen) {
		once.Do(func() {
			screen.Clear()
			refresh(screen)
			screen.Show()
		})
	}
}

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
	mutate dscope.Mutate,
	extraAttrs ShowExtraAttrs,
	attrsBox AttrsBox,
) DashboardRefresh {
	return func(screen tcell.Screen) {

		width, height := screen.Size()
		var handlers ClickHandlers

		// evs
		x := 0
		y := height - 1
		if offset < 0 {
			y += int(offset)
			offset = 0
		}
		for i := len(evs) - 1 - int(offset); i >= 0; i-- {
			ev := evs[i]
			box := evBox(ev)
			dx, dy := box.DrawLeftBottom(screen, x, y, width, height)
			y -= dy
			if y < 0 {
				break
			}
			handlers = append(handlers, ClickHandler{
				X:      x,
				Y:      y,
				Width:  dx,
				Height: dy,
				Func: func(
					mutate dscope.Mutate,
				) {
					mutate(func() ShowExtraAttrs {
						return ev.ExtraAttrs
					})
				},
			})
		}

		// extra attrs
		if len(extraAttrs) > 0 {
			box := attrsBox(extraAttrs)
			box.Draw(
				screen,
				width/2, 0,
				width, height,
			)
		}

		mutate(&handlers)

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

type ClickHandler struct {
	X, Y          int
	Width, Height int
	Func          any
}

type ClickHandlers []ClickHandler

func (_ DashboardDef) ClickHandlers() ClickHandlers {
	return nil
}

type ShowExtraAttrs []Attr

func (_ DashboardDef) ShowExtraAttrs() ShowExtraAttrs {
	return nil
}

type AttrsBox func(attrs []Attr) BoxSpec

func (_ DashboardDef) AttrsBox(
	styles Styles,
) AttrsBox {
	return func(attrs []Attr) BoxSpec {
		var rows []RowSpec
		for _, attr := range attrs {
			var row RowSpec
			row.AppendStr(attr.Name, styles.AttrName)
			row.AppendStr(": ", tcell.StyleDefault)
			row.AppendStr(fmt.Sprintf("%v", attr.Value), tcell.StyleDefault)
			rows = append(rows, row)
		}
		return BoxSpec{
			Rows: rows,
		}
	}
}
