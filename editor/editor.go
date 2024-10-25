package editor

import (
	"bytes"
	"log"
	"simpledit/buffer"
	"simpledit/screen"

	"github.com/gdamore/tcell/v2"
)

type Screen = tcell.Screen

type Editor struct {
	screen     screen.EditorScreen
	buffer     []byte
	bufferRows [][]byte
}

func NewEditor() (*Editor, error) {
	editorScreen, err := screen.InitEditorScreen()

	if err != nil {
		log.Fatalln(err)
	}

	editor := Editor{
		screen: *editorScreen,
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
