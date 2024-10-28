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
	buffer     []byte
	bufferRows [][]byte
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
		buffer: []byte{},
	}

	return &editor, nil
}

func (editor *Editor) ReadFileIntoBuffer(fileName string) {
	editor.buffer = buffer.ReadFile(fileName)
	editor.bufferRows = bytes.Split(editor.buffer, []byte{'\n'})
}

func (editor *Editor) GetBufferRows() [][]byte {
	return editor.bufferRows
}

func (editor *Editor) GetScreen() Screen {
	return editor.screen.GetScreen()
}

func (editor *Editor) Render() {
	editor.screen.DrawBufferRows(editor.bufferRows)
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

			if int(c.Col) > len(editor.GetBufferRows()[c.Row]) && int(c.Row) < len(editor.GetBufferRows()) {
				c.Row++
				c.Col = 0
			}

			break
		}
		if event.Key() == tcell.KeyLeft {
			c.Col--

			if c.Col < 0 {
				if c.Row > 0 {
					c.Row--
					c.Col = len(editor.GetBufferRows()[c.Row])
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
			// bufferRows = slices.Insert(bufferRows, cursorRow, []byte{'p', 'l', 's'})
			// cursorRow++
			// cursorCol = 0
			break
		}
		if event.Key() == tcell.KeyTab {
			editor.GetBufferRows()[c.Row] = slices.Insert(editor.GetBufferRows()[c.Row], int(c.Col), byte('\t'))
			editor.GetBufferRows()[c.Row] = slices.Insert(editor.GetBufferRows()[c.Row], int(c.Col), byte('\t'))
			editor.GetBufferRows()[c.Row] = slices.Insert(editor.GetBufferRows()[c.Row], int(c.Col), byte('\t'))
			editor.GetBufferRows()[c.Row] = slices.Insert(editor.GetBufferRows()[c.Row], int(c.Col), byte('\t'))
			c.Col += 4
			break
		}
		if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			if c.Col <= 0 {
				break
			}

			row := editor.GetBufferRows()[c.Row]

			if c.Col < len(row) {
				editor.GetBufferRows()[c.Row] = slices.Delete(editor.GetBufferRows()[c.Row], c.Col, c.Col+1)

				// editor.GetBufferRows()[c.Row] = []byte{}
				// editor.GetBufferRows()[c.Row] = append(row[:c.Col], row[c.Col+1:]...)
			} else {
				editor.GetBufferRows()[c.Row] = row[:c.Col]
			}

			break
		}

		char := event.Rune()
		editor.GetBufferRows()[c.Row] = slices.Insert(editor.GetBufferRows()[c.Row], int(c.Col), byte(char))
		c.Col++
	}

	// Safeguarding against overflow
	if c.Row < 0 {
		c.Row = 0
	}

	if int(c.Row) >= len(editor.GetBufferRows()) {
		c.Row = len(editor.GetBufferRows()) - 1
	}

	if c.Col < 0 {
		c.Col = 0
	}

	if int(c.Col) > len(editor.GetBufferRows()[c.Row]) {
		c.Col = max(len(editor.GetBufferRows()[c.Row])-1, 0)
	}
}

func (editor *Editor) ShowCursor() {
	editor.GetScreen().ShowCursor(int(editor.cursor.Col), int(editor.cursor.Row))
}
