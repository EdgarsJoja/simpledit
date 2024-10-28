package main

import (
	"log"
	"simpledit/editor"
)

func main() {
	editor, err := editor.NewEditor()

	if err != nil {
		log.Fatalln(err)
	}

	editor.ReadFileIntoBuffer("test.txt")

	for {
		editor.Render()
		editor.HandleKeyEvents()
		editor.ShowCursor()
	}
}
