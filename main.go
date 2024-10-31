package main

import (
	"log"
	"os"
	"simpledit/editor"
)

func main() {
	if len(os.Args) < 2 {
		panic("File not specified")
	}

	fileName := os.Args[1]

	editor, err := editor.NewEditor()

	if err != nil {
		log.Fatalln(err)
	}

	editor.ReadFileIntoBuffer(fileName)

	for {
		editor.Render()
		editor.HandleEvents()
		editor.ShowCursor()
	}
}
