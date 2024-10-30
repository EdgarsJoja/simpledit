package editor

import (
	"bytes"
	"log"
	"os"
	"simpledit/buffer"
	"simpledit/screen"
	"slices"

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
	editor.screen.GetScreen().Clear()
	editor.screen.DrawBufferRows(editor.BufferRows)
	editor.GetScreen().Show()
}

func (editor *Editor) HandleKeyEvents() {
	event := editor.GetScreen().PollEvent()
	c := &editor.cursor

	switch event := event.(type) {
	case *tcell.EventKey:
		if event.Key() == tcell.KeyEscape {
			os.Exit(0)

			return
		}
		if event.Key() == tcell.KeyRight {
			c.Col++

			if int(c.Col) > len(editor.GetCurrentRow()) && int(c.Row) < len(editor.BufferRows)-1 {
				editor.CursorGoToStartOfNextRow()
			}

			break
		}
		if event.Key() == tcell.KeyLeft {
			c.Col--

			if c.Col < 0 {
				if c.Row > 0 {
					editor.CursorGoToEndOfPreviousRow()
				} else {
					c.Col = 0
				}
			}

			break
		}
		if event.Key() == tcell.KeyUp {
			c.Row--
			break
		}
		if event.Key() == tcell.KeyDown {
			c.Row++
			break
		}
		if event.Key() == tcell.KeyEnter {
			// Insert new line above current line
			if c.Col == 0 {
				editor.BufferRows = slices.Insert(editor.BufferRows, c.Row, []byte{})
				c.Row++
				break
			}

			// Insert new line below current line
			if c.Col >= len(editor.GetCurrentRow())-1 {
				editor.BufferRows = slices.Insert(editor.BufferRows, c.Row+1, []byte{})
				c.Row++
				break
			}

			// Split current row
			row := editor.GetCurrentRow()
			editor.SetCurrentRow(row[:c.Col])
			editor.BufferRows = slices.Insert(editor.BufferRows, c.Row+1, row[c.Col:])

			c.Row++
			c.Col = 0
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

			c.Col--

			break
		}

		if event.Key() == tcell.KeyCtrlS {
			editor.WriteBufferToFile()
			os.Exit(0)
			break
		}

		char := event.Rune()
		editor.SetCurrentRow(slices.Insert(editor.GetCurrentRow(), int(c.Col), byte(char)))
		c.Col++
	}

	// Safeguarding against overflow
	if c.Row < 0 {
		c.Row = 0
	}

	if int(c.Row) >= len(editor.BufferRows) {
		c.Row = len(editor.BufferRows) - 1
	}

	if c.Col < 0 {
		c.Col = 0
	}

	if int(c.Col) > len(editor.GetCurrentRow()) {
		c.Col = max(len(editor.GetCurrentRow()), 0)
	}
}

func (editor *Editor) ShowCursor() {
	editor.GetScreen().ShowCursor(int(editor.cursor.Col), int(editor.cursor.Row))
}

func (editor *Editor) CursorGoToEndOfPreviousRow() {
	editor.cursor.Row--
	editor.cursor.Col = len(editor.GetCurrentRow())
}

func (editor *Editor) CursorGoToStartOfNextRow() {
	editor.cursor.Row++
	editor.cursor.Col = 0
}
