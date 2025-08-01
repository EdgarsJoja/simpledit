package screen

import (
	"github.com/gdamore/tcell/v2"
)

type Coordinates struct {
	Row int
	Col int
}

type EditorScreen struct {
	screen         tcell.Screen
	StartRow       int
	EndRow         int
	StartCol       int
	EndCol         int
	ScreenWidth    int
	ScreenHeight   int
	HighlightStart *Coordinates
	HighlightEnd   *Coordinates
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

func (editorScreen *EditorScreen) HighlightStyle() tcell.Style {
	return tcell.StyleDefault.Background(tcell.ColorDimGray).Foreground(tcell.ColorReset)
}

func (editorScreen *EditorScreen) DrawRune(col int, row int, value rune, style tcell.Style) {
	editorScreen.screen.SetContent(col, row, value, nil, style)
}

func (editorScreen *EditorScreen) DrawText(col int, row int, value string) {
	x := col

	for _, r := range value {
		editorScreen.DrawRune(x, row, r, editorScreen.DefaultStyle())
		x++
	}
}

func (editorScreen *EditorScreen) DrawBufferRows(bufferRows [][]byte) {
	rowMultiplier := 1_000_000

	for row, bufferRow := range bufferRows {
		for col, char := range bufferRow {
			style := editorScreen.DefaultStyle()

			hStart := editorScreen.HighlightStart
			hEnd := editorScreen.HighlightEnd
			highlight := hStart != nil && hEnd != nil

			if highlight {
				pos := row*rowMultiplier + col
				hStartPos := hStart.Row*rowMultiplier + hStart.Col
				hEndPos := hEnd.Row*rowMultiplier + hEnd.Col

				inHighlightRange := highlight && pos >= min(hStartPos, hEndPos) && pos <= max(hStartPos, hEndPos)

				if inHighlightRange {
					style = editorScreen.HighlightStyle()
				}
			}

			editorScreen.DrawRune(col, row, rune(char), style)
		}
	}
}
