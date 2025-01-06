package ppt

import "context"

type NotesPersistAtom struct {
	Record
	persistIdRef   uint32
	notesContainer *NotesContainer
}

func (c *NotesPersistAtom) Parse(ctx context.Context, persistIdOffset map[uint32]int64) {
	c.persistIdRef = c.Uint32At(0)
	notes, err := readRecord()
}
