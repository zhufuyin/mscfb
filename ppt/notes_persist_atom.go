package ppt

import "context"

type NotesPersistAtom struct {
	Record
	persistIdRef uint32
}

func (c *NotesPersistAtom) parse(ctx context.Context) error {
	c.persistIdRef = c.Uint32At(0)
	return nil
}
