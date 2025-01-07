package ppt

import (
	"context"
	"fmt"
)

type DocumentContainer struct {
	Record
	masterList *MasterListWithTextContainer
	slideList  *SlideListWithTextContainer
	notesList  *NotesListWithTextContainer
}

func (c *DocumentContainer) parse(ctx context.Context) error {
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
				slideList := &SlideListWithTextContainer{
					Record: record,
				}
				err = slideList.parse()
				if err != nil {
					return err
				}
				c.slideList = slideList
			case 0x001: // MasterListWithTextContainer
				c.masterList = &MasterListWithTextContainer{
					Record: record,
				}
			case 0x002: // NotesListWithTextContainer
				notesList := &NotesListWithTextContainer{
					Record: record,
				}
				err = notesList.parse(ctx)
				if err != nil {
					return err
				}
				c.notesList = notesList
			}
			//if c.masterList != nil && c.slideList != nil && c.notesList != nil {
			//	return nil
			//}
		}
		offset += record.DataLength + headerSize
	}
	return nil
}

func (c *DocumentContainer) extractText() ([]string, error) {
	var texts []string
	if c.slideList == nil {
		return nil, nil
	}
	txts, err := c.slideList.extractText()
	if err != nil {
		return nil, err
	}
	if len(txts) > 0 {
		texts = append(texts, txts...)
	}
	// todo notes
	return texts, nil
}
