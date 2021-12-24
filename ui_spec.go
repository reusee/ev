package ev

import "github.com/gdamore/tcell/v2"

type CellSpec struct {
	Rune  rune
	Style tcell.Style
}

func (c CellSpec) Draw(screen tcell.Screen, x, y, width, height int) int {
	screen.SetContent(x, y, c.Rune, nil, c.Style)
	return runeDisplayWidth(c.Rune)
}

type RowSpec struct {
	Cells []CellSpec
}

func (r *RowSpec) AppendStr(str string, style tcell.Style) {
	for _, ru := range str {
		r.Cells = append(r.Cells, CellSpec{
			Rune:  ru,
			Style: style,
		})
	}
}

func (r RowSpec) Draw(screen tcell.Screen, x, y, width, height int) int {
	originX := x
	for _, cell := range r.Cells {
		x += cell.Draw(screen, x, y, width, height)
		if x >= width {
			return x - originX - 1
		}
	}
	return x - originX
}

type BoxSpec struct {
	Rows []RowSpec
}

func (b BoxSpec) Draw(screen tcell.Screen, x, y, width, height int) int {
	originY := y
	for _, row := range b.Rows {
		row.Draw(screen, x, y, width, height)
		if y >= height {
			return y - originY
		}
		y++
	}
	return y - originY
}

func (b BoxSpec) DrawLeftBottom(screen tcell.Screen, x, y, width, height int) int {
	originY := y
	for i := len(b.Rows) - 1; i >= 0; i-- {
		row := b.Rows[i]
		row.Draw(screen, x, y, width, height)
		y--
		if y < 0 {
			return y - originY
		}
	}
	return y - originY
}
