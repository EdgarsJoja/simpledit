package main

import (
	"log"
	"simpledit/editor"
	"slices"

	"github.com/gdamore/tcell/v2"
)

func main() {
	editor, err := editor.NewEditor()

	if err != nil {
		log.Fatalln(err)
	}

	editor.ReadFileIntoBuffer("test.txt")
	bufferRows := editor.GetBufferRows()

	cursorCol, cursorRow := 0, 0

	for {
		editor.Render()

		event := editor.GetScreen().PollEvent()

		switch event := event.(type) {
		case *tcell.EventKey:
			if event.Key() == tcell.KeyEscape {
				return
			}
			if event.Key() == tcell.KeyRight {
				cursorCol++

				if cursorCol > len(editor.GetBufferRows()[cursorRow]) && cursorRow < len(editor.GetBufferRows()) {
					cursorRow++
					cursorCol = 0
				}

				break
			}
			if event.Key() == tcell.KeyLeft {
				cursorCol--

				if cursorCol < 0 {
					if cursorRow > 0 {
						cursorRow--
						cursorCol = len(editor.GetBufferRows()[cursorRow])
					} else {
						cursorCol = 0
					}
				}

				break
			}
			if event.Key() == tcell.KeyUp {
				cursorRow--
				break
			}
			if event.Key() == tcell.KeyDown {
				cursorRow++
				break
			}
			if event.Key() == tcell.KeyEnter {
				// bufferRows = slices.Insert(bufferRows, cursorRow, []byte{'p', 'l', 's'})
				// cursorRow++
				// cursorCol = 0
				break
			}
			if event.Key() == tcell.KeyTab {
				bufferRows[cursorRow] = slices.Insert(bufferRows[cursorRow], cursorCol, byte('\t'))
				bufferRows[cursorRow] = slices.Insert(bufferRows[cursorRow], cursorCol, byte('\t'))
				bufferRows[cursorRow] = slices.Insert(bufferRows[cursorRow], cursorCol, byte('\t'))
				bufferRows[cursorRow] = slices.Insert(bufferRows[cursorRow], cursorCol, byte('\t'))
				cursorCol += 4
				break
			}

			char := event.Rune()
			editor.GetBufferRows()[cursorRow] = slices.Insert(editor.GetBufferRows()[cursorRow], cursorCol, byte(char))
			cursorCol++
		}

		// Safeguarding agains overflow
		if cursorRow < 0 {
			cursorRow = 0
		}

		if cursorRow >= len(editor.GetBufferRows()) {
			cursorRow = len(editor.GetBufferRows()) - 1
		}

		if cursorCol < 0 {
			cursorCol = 0
		}

		if cursorCol > len(editor.GetBufferRows()[cursorRow]) {
			cursorCol = max(len(editor.GetBufferRows()[cursorRow])-1, 0)
		}

		editor.GetScreen().ShowCursor(cursorCol, cursorRow)
	}
}
