package repocounter

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

func progressBar(screen tcell.Screen, width int, percentage float64, x, y int) {
	progressSegments := int(float64(width) * percentage)
	negativeSegments := width - progressSegments

	for i := 0; i < progressSegments; i++ {
		screen.SetCell(x+i, y, tcell.StyleDefault.Background(tcell.ColorAntiqueWhite), ' ')
	}
	for i := 0; i < negativeSegments; i++ {
		screen.SetCell(x+progressSegments+i, y, tcell.StyleDefault.Background(tcell.ColorDimGrey), ' ')
	}
}

func progressBarWithCounter(screen tcell.Screen, width int, total, parts, x, y int) {
	totalRepoLength := log10(total) // Number of digits
	counter := fmt.Sprintf("%0*d/%d ", totalRepoLength, parts, total)
	percentage := float64(parts) / float64(total)
	emitStr(screen, x, y, tcell.StyleDefault, counter)

	progressBar(screen, width-len(counter), percentage, x+len(counter), y)
}

func button(screen tcell.Screen, x, y int, text string, selected bool) int {
	style := tcell.StyleDefault
	if selected {
		style = tcell.StyleDefault.Background(tcell.ColorDimGray)
	}
	firstStyle := style.Underline(true)

	emitStr(screen, x, y, style, "[ ")
	emitStr(screen, x+2, y, firstStyle, text[:1])
	emitStr(screen, x+3, y, style, text[1:]+" ]")

	return len(text) + 4
}

func center(str string, size int) string {
	spaces := size - len(str)
	before := spaces / 2
	after := spaces - before

	strBuilder := strings.Builder{}
	for i := 0; i < before; i++ {
		strBuilder.WriteString(" ")
	}
	strBuilder.WriteString(str)
	for i := 0; i < after; i++ {
		strBuilder.WriteString(" ")
	}
	return strBuilder.String()
}

func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) int {
	width := 0
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		width += w
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
	return width
}
