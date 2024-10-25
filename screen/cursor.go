package screen

type Cursor struct {
	row uint
	col uint
}

func (cursor *Cursor) goRight() {
	cursor.col++
}
