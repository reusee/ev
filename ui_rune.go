package ev

import (
	textwidth "golang.org/x/text/width"
	"sync"
)

var runeWidths sync.Map

func runeDisplayWidth(r rune) int {
	if v, ok := runeWidths.Load(r); ok {
		return v.(int)
	}
	prop := textwidth.LookupRune(r)
	kind := prop.Kind()
	width := 1
	if kind == textwidth.EastAsianAmbiguous ||
		kind == textwidth.EastAsianWide ||
		kind == textwidth.EastAsianFullwidth {
		width = 2
	}
	runeWidths.Store(r, width)
	return width
}
