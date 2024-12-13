package screen

type Cursor struct {
	Row             int
	Col             int
	updatedCallback func(cursor *Cursor)
}

func (cursor *Cursor) SetCol(value int) {
	cursor.Col = value
	cursor.updatedCallback(cursor)
}

func (cursor *Cursor) SetRow(value int) {
	cursor.Row = value
	cursor.updatedCallback(cursor)
}

func (cursor *Cursor) SetUpdatedCallback(callback func(cursor *Cursor)) {
	cursor.updatedCallback = callback
}
