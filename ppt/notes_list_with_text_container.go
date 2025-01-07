package ppt

import "context"

type NotesListWithTextContainer struct {
	Record
	rgNotesPersistAtom []*NotesPersistAtom
}

func (c *NotesListWithTextContainer) parse(ctx context.Context) error {
	offset := int64(0)
	for offset < c.DataLength {
		record, err := readRecord(c, offset, recordTypeSlidePersistAtom)
		if err != nil {
			return err
		}
		notesPersistAtom := &NotesPersistAtom{
			Record: record,
		}
		err = notesPersistAtom.parse(ctx)
		if err != nil {
			return err
		}
		c.rgNotesPersistAtom = append(c.rgNotesPersistAtom, notesPersistAtom)
		offset += record.DataLength + headerSize
	}
	return nil
}
