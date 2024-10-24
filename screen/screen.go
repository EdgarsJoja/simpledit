package screen

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

func InitScreen() tcell.Screen {
	screen, err := tcell.NewScreen()

	if err != nil {
		log.Fatalln(err)
	}

	if err := screen.Init(); err != nil {
		log.Fatalln(err)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	screen.SetStyle(defStyle)

	screen.Clear()

	return screen
}

func DrawBuffer(tScreen tcell.Screen, buffer [][]byte) {
	for row, bufferRow := range buffer {
		for col, char := range bufferRow {
			tScreen.SetContent(col, row, rune(char), nil, DefaultStyle())
		}
	}
}

func DefaultStyle() tcell.Style {
	return tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
}
