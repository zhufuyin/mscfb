package ppt

type NotesListWithTextContainer struct {
	Record
	rgNotesPersistAtom []*NotesPersistAtom
}

func (c *NotesListWithTextContainer) parse() error {
	offset := int64(0)
	for offset < c.DataLength {
		record, err := readRecord(c, offset, recordTypeUnspecified)
		if err != nil {
			return err
		}
		notesPersistAtom := &NotesPersistAtom{
			Record: record,
		}
		c.rgNotesPersistAtom = append(c.rgNotesPersistAtom, notesPersistAtom)
		offset += record.DataLength + headerSize
	}
	return nil
}
