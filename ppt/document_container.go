package ppt

import "fmt"

type DocumentContainer struct {
	Record
	masterList *MasterListWithTextContainer
	slideList  *SlideListWithTextContainer
	notesList  *NotesListWithTextContainer
}

func (c *DocumentContainer) parse() error {
	if c.RecType != recordTypeDocument {
		return fmt.Errorf("invalid document record type %d", c.RecType)
	}
	// skip documentAtom elem bytes
	offset := int64(48)
	for offset < c.DataLength {
		record, err := readRecordHeaderOnly(c, offset, recordTypeUnspecified)
		if err != nil {
			return err
		}
		if record.RecType == recordTypeSlideListWithText {
			record.RecordData = make([]byte, record.DataLength)
			_, err := c.ReadAt(record.RecordData, offset+headerSize)
			if err != nil {
				return err
			}
			switch record.TypeInstance {
			case 0x000: // SlideListWithTextContainer
				c.slideList = &SlideListWithTextContainer{
					Record: record,
				}
			case 0x001: // MasterListWithTextContainer
				c.masterList = &MasterListWithTextContainer{
					Record: record,
				}
			case 0x002: // NotesListWithTextContainer
				c.notesList = &NotesListWithTextContainer{
					Record: record,
				}
			}
			if c.masterList != nil && c.slideList != nil && c.notesList != nil {
				return nil
			}
		}
		offset += record.DataLength + headerSize
	}
	return nil
}
