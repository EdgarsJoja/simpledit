package screen

import (
	"github.com/gdamore/tcell/v2"
)

type EditorScreen struct {
	screen       tcell.Screen
	StartRow     int
	EndRow       int
	StartCol     int
	EndCol       int
	ScreenWidth  int
	ScreenHeight int
}

func InitEditorScreen() (*EditorScreen, error) {
	tScreen, err := tcell.NewScreen()

	if err != nil {
		return nil, err
	}

	if err := tScreen.Init(); err != nil {
		return nil, err
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	tScreen.SetStyle(defStyle)

	tScreen.Clear()

	editorScreen := EditorScreen{
		screen: tScreen,
	}

	return &editorScreen, nil
}

func (editorScreen *EditorScreen) GetScreen() tcell.Screen {
	return editorScreen.screen
}

func (editorScreen *EditorScreen) DefaultStyle() tcell.Style {
	return tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
}

func (editorScreen *EditorScreen) DrawRune(col int, row int, value rune) {
	editorScreen.screen.SetContent(col, row, value, nil, editorScreen.DefaultStyle())
}

func (editorScreen *EditorScreen) DrawText(col int, row int, value string) {
	x := col

	for _, r := range value {
		editorScreen.DrawRune(x, row, r)
		x++
	}
}

func (editorScreen *EditorScreen) DrawBufferRows(bufferRows [][]byte) {
	for row, bufferRow := range bufferRows {
		for col, char := range bufferRow {
			editorScreen.DrawRune(col, row, rune(char))
		}
	}
}
