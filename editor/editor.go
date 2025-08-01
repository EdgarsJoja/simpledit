package editor

import (
	"bytes"
	"log"
	"os"
	"simpledit/buffer"
	"simpledit/screen"
	"slices"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type Screen = tcell.Screen

type Editor struct {
	screen     screen.EditorScreen
	cursor     screen.Cursor
	BufferRows [][]byte
	fileName   string
}

func NewEditor() (*Editor, error) {
	editorScreen, err := screen.InitEditorScreen()

	if err != nil {
		log.Fatalln(err)
	}

	editorCursor := &screen.Cursor{
		Row: 0,
		Col: 0,
	}

	editor := Editor{
		screen: *editorScreen,
		cursor: *editorCursor,
	}

	editor.cursor.SetUpdatedCallback(func(cursor *screen.Cursor) {
		if cursor.Col < editor.screen.StartCol {
			editor.screen.EndCol = editor.cursor.Col + editor.screen.ScreenWidth + 1
			editor.screen.StartCol = max(editor.screen.EndCol-editor.screen.ScreenWidth-1, 0)
		}

		if cursor.Col >= editor.screen.EndCol {
			editor.screen.StartCol = max(editor.cursor.Col-editor.screen.ScreenWidth+1, 0)
			editor.screen.EndCol = editor.screen.StartCol + editor.screen.ScreenWidth
		}

		if cursor.Row < editor.screen.StartRow {
			editor.screen.EndRow = editor.cursor.Row + editor.screen.ScreenHeight + 1
			editor.screen.StartRow = max(editor.screen.EndRow-editor.screen.ScreenHeight-1, 0)
		}

		if cursor.Row >= editor.screen.EndRow {
			editor.screen.StartRow = max(editor.cursor.Row-editor.screen.ScreenHeight+1, 0)
			editor.screen.EndRow = editor.screen.StartRow + editor.screen.ScreenHeight
		}
	})

	width, height := editor.GetScreen().Size()
	editor.screen.ScreenHeight = height
	editor.screen.ScreenWidth = width

	return &editor, nil
}

func (editor *Editor) ReadFileIntoBuffer(fileName string) {
	editor.fileName = fileName

	buffer := buffer.ReadFile(fileName)
	editor.BufferRows = bytes.Split(buffer, []byte{'\n'})
}

func (editor *Editor) WriteBufferToFile() {
	buffer.WriteFile(editor.fileName, editor.BufferRows)
}

func (editor *Editor) GetCurrentRow() []byte {
	return editor.BufferRows[editor.cursor.Row]
}

func (editor *Editor) SetCurrentRow(row []byte) {
	editor.BufferRows[editor.cursor.Row] = row
}

func (editor *Editor) GetScreen() Screen {
	return editor.screen.GetScreen()
}

func (editor *Editor) Render() {
	editor.GetScreen().Clear()

	bufferRows := make([][]byte, len(editor.BufferRows))

	rowLowerBound, rowUpperBound := max(editor.screen.StartRow, 0), min(editor.screen.EndRow, len(editor.BufferRows))
	_ = copy(bufferRows, editor.BufferRows[rowLowerBound:rowUpperBound])

	for i, row := range bufferRows {
		colLowerBound := max(min(editor.screen.StartCol, len(row)), 0)
		colUpperBound := min(editor.screen.EndCol, max(len(row), 0))
		bufferRows[i] = bufferRows[i][colLowerBound:colUpperBound]
	}

	editor.screen.DrawBufferRows(bufferRows)

	colRow := strconv.Itoa(editor.cursor.Row) + ":" + strconv.Itoa(editor.cursor.Col)
	editor.screen.DrawText(0, editor.screen.EndRow-1, colRow)

	editor.GetScreen().Show()
}

func (editor *Editor) HandleEvents() {
	event := editor.GetScreen().PollEvent()
	c := &editor.cursor
	s := &editor.screen

	switch event := event.(type) {
	case *tcell.EventKey:
		shift := strings.Contains(event.Name(), "Shift")
		if shift {
			cursorCoordinates := screen.Coordinates{Row: c.Row, Col: c.Col}

			if s.HighlightStart == nil || s.HighlightEnd == nil {
				s.HighlightStart = &cursorCoordinates
				s.HighlightEnd = &cursorCoordinates
			} else {
				s.HighlightEnd = &cursorCoordinates
			}
		} else {
			s.HighlightStart = nil
			s.HighlightEnd = nil
		}
		if event.Key() == tcell.KeyEscape {
			os.Exit(0)

			return
		}
		if event.Key() == tcell.KeyRight {
			c.SetCol(c.Col + 1)

			if int(c.Col) > len(editor.GetCurrentRow()) && int(c.Row) < len(editor.BufferRows)-1 {
				editor.CursorGoToStartOfNextRow()
			}

			c.TargetCol = c.Col

			break
		}
		if event.Key() == tcell.KeyLeft {
			c.SetCol(c.Col - 1)

			if c.Col < 0 {
				if c.Row > 0 {
					editor.CursorGoToEndOfPreviousRow()
				} else {
					c.SetCol(0)
				}
			}

			c.TargetCol = c.Col

			break
		}
		if event.Key() == tcell.KeyUp {
			c.SetRow(c.Row - 1)
			c.SetCol(c.TargetCol)
			break
		}
		if event.Key() == tcell.KeyDown {
			c.SetRow(c.Row + 1)
			c.SetCol(c.TargetCol)
			break
		}
		if event.Key() == tcell.KeyEnter {
			// Insert new line above current line
			if c.Col == 0 {
				editor.BufferRows = slices.Insert(editor.BufferRows, c.Row, []byte{})
				c.SetRow(c.Row + 1)
				break
			}

			// Insert new line below current line
			if c.Col >= len(editor.GetCurrentRow())-1 {
				editor.BufferRows = slices.Insert(editor.BufferRows, c.Row+1, []byte{})
				c.SetRow(c.Row + 1)
				break
			}

			// Split current row
			row := editor.GetCurrentRow()
			editor.SetCurrentRow(row[:c.Col])
			editor.BufferRows = slices.Insert(editor.BufferRows, c.Row+1, row[c.Col:])

			c.SetRow(c.Row + 1)
			c.SetCol(0)
			break
		}
		if event.Key() == tcell.KeyTab {
			// editor.SetCurrentRow(slices.Insert(editor.GetCurrentRow(), int(c.Col), byte('\t')))
			// editor.SetCurrentRow(slices.Insert(editor.GetCurrentRow(), int(c.Col), byte('\t')))
			// editor.SetCurrentRow(slices.Insert(editor.GetCurrentRow(), int(c.Col), byte('\t')))
			// editor.SetCurrentRow(slices.Insert(editor.GetCurrentRow(), int(c.Col), byte('\t')))
			// c.Col += 4
			break
		}
		if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			if c.Col == 0 && c.Row <= 0 {
				break
			}

			// Remove current row, if it was empty
			if len(editor.GetCurrentRow()) == 0 {
				editor.BufferRows = slices.Delete(editor.BufferRows, c.Row, c.Row+1)
				editor.CursorGoToEndOfPreviousRow()
				break
			}

			// Go to the end of previous row
			if c.Col == 0 {
				row := editor.GetCurrentRow()
				editor.BufferRows = slices.Delete(editor.BufferRows, c.Row, c.Row+1)
				editor.CursorGoToEndOfPreviousRow()
				editor.SetCurrentRow(append(editor.GetCurrentRow(), row...))
				break
			}

			// Delete character
			if c.Col <= len(editor.GetCurrentRow())-1 {
				editor.SetCurrentRow(slices.Delete(editor.GetCurrentRow(), c.Col-1, c.Col))
			} else {
				editor.SetCurrentRow(editor.GetCurrentRow()[:c.Col-1])
			}

			c.SetCol(c.Col - 1)

			break
		}

		if event.Key() == tcell.KeyCtrlS {
			editor.WriteBufferToFile()
			os.Exit(0)
			break
		}

		char := event.Rune()
		editor.SetCurrentRow(slices.Insert(editor.GetCurrentRow(), int(c.Col), byte(char)))
		c.SetCol(c.Col + 1)
		c.TargetCol = c.Col
	case *tcell.EventResize:
		_, height := event.Size()
		editor.screen.StartRow = 0
		editor.screen.EndRow = height

		c.SetRow(editor.screen.StartRow)
	}

	// Safeguarding against overflow
	if c.Row < 0 {
		c.SetRow(0)
	}

	if c.Row >= len(editor.BufferRows) {
		c.SetRow(len(editor.BufferRows) - 1)
	}

	if c.Col < 0 {
		c.SetCol(0)
	}

	if c.Col > len(editor.GetCurrentRow()) {
		c.SetCol(max(len(editor.GetCurrentRow()), 0))
	}
}

func (editor *Editor) ShowCursor() {
	editor.GetScreen().ShowCursor(editor.cursor.Col-editor.screen.StartCol, editor.cursor.Row-editor.screen.StartRow)
}

func (editor *Editor) CursorGoToEndOfPreviousRow() {
	editor.cursor.SetRow(editor.cursor.Row - 1)
	editor.cursor.SetCol(len(editor.GetCurrentRow()))
}

func (editor *Editor) CursorGoToStartOfNextRow() {
	editor.cursor.SetRow(editor.cursor.Row + 1)
	editor.cursor.SetCol(0)
}
