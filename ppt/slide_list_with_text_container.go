package ppt

// refer MS-PPT 2.4.14.6
type SlideListWithTextContainer struct {
	Record
	slidePersistAtoms []*SlidePersistAtom
	textCharsAtoms    []*TextCharsAtom
	textBytesAtoms    []*TextBytesAtom
}

func (c *SlideListWithTextContainer) parse() error {
	offset := int64(0)
	for offset < c.DataLength {
		record, err := readRecordHeaderOnly(c, offset, recordTypeUnspecified)
		if err != nil {
			return err
		}
		switch record.RecType {
		case recordTypeSlidePersistAtom:
			err := readRecordData(c, &record, offset)
			if err != nil {
				return err
			}
			c.slidePersistAtoms = append(c.slidePersistAtoms, &SlidePersistAtom{
				Record: record,
			})
		case recordTypeTextCharsAtom:
			err := readRecordData(c, &record, offset)
			if err != nil {
				return err
			}
			textCharsAtom := TextCharsAtom(record)
			c.textCharsAtoms = append(c.textCharsAtoms, &textCharsAtom)
		case recordTypeTextBytesAtom:
			err := readRecordData(c, &record, offset)
			if err != nil {
				return err
			}
			textBytesAtom := TextBytesAtom(record)
			c.textBytesAtoms = append(c.textBytesAtoms, &textBytesAtom)
		}
		offset += record.DataLength + headerSize
	}
	return nil
}
