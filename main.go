package main

import (
	"bytes"
	"simpledit/buffer"
	"simpledit/screen"
	"slices"

	"github.com/gdamore/tcell/v2"
)

func main() {
	buffer := buffer.ReadFile("test.txt")
	tScreen := screen.InitScreen()

	bufferRows := bytes.Split(buffer, []byte{'\n'})

	cursorCol, cursorRow := 0, 0

	for {
		screen.DrawBuffer(tScreen, bufferRows)
		tScreen.Show()

		event := tScreen.PollEvent()

		switch event := event.(type) {
		case *tcell.EventKey:
			if event.Key() == tcell.KeyEscape {
				return
			}
			if event.Key() == tcell.KeyRight {
				cursorCol++

				if cursorCol > len(bufferRows[cursorRow]) && cursorRow < len(bufferRows) {
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
						cursorCol = len(bufferRows[cursorRow])
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
			bufferRows[cursorRow] = slices.Insert(bufferRows[cursorRow], cursorCol, byte(char))
			cursorCol++
		}

		// Safeguarding agains overflow
		if cursorRow < 0 {
			cursorRow = 0
		}

		if cursorRow >= len(bufferRows) {
			cursorRow = len(bufferRows) - 1
		}

		if cursorCol < 0 {
			cursorCol = 0
		}

		if cursorCol > len(bufferRows[cursorRow]) {
			cursorCol = max(len(bufferRows[cursorRow])-1, 0)
		}

		tScreen.ShowCursor(cursorCol, cursorRow)
	}
}
