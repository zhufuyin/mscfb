package ppt

import (
	"context"
	"fmt"
)

type DocumentContainer struct {
	Record
	pptDocStream *PowerPointStream
	masterList   *MasterListWithTextContainer
	slideList    *SlideListWithTextContainer
	notesList    *NotesListWithTextContainer
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
				masterList := &MasterListWithTextContainer{
					Record: record,
				}
				err = masterList.parse(ctx)
				if err != nil {
					return err
				}
				c.masterList = masterList
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
	pptDocStream := c.pptDocStream
	// range items in SlideListWithTextContainer
	for _, item := range c.slideList.items {
		switch atom := item.(type) {
		case *SlidePersistAtom:
			txts, err := pptDocStream.extractTextFromSlide(atom)
			if err != nil {
				return nil, err
			}
			texts = append(texts, txts...)
		case *TextCharsAtom:
			txt, err := atom.getText()
			if err != nil {
				return nil, err
			}
			texts = append(texts, txt)
		case *TextBytesAtom:
			txt, err := atom.getText()
			if err != nil {
				return nil, err
			}
			texts = append(texts, txt)
		default:
			fmt.Printf("unknown item type in SlideListWithTextContainer: %t", item)
		}
	}

	//txts, err := c.slideList.extractText()
	//if err != nil {
	//	return nil, err
	//}
	//if len(txts) > 0 {
	//	texts = append(texts, txts...)
	//}

	return texts, nil
}
