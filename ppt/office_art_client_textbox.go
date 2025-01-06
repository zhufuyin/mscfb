package ppt

import "fmt"

type OfficeArtClientTextbox struct {
	Record
	texts []string
}

func (o *OfficeArtClientTextbox) readTextAtoms() error {
	offset := int64(0)
	var text string
	for offset < o.DataLength {
		rec, err := readRecordHeaderOnly(o, offset, recordTypeUnspecified)
		if err != nil {
			return err
		}
		//fmt.Printf("record type=%X instance=%X\n", rec.Type(), rec.Instance())
		if rec.Type() == recordTypeTextCharsAtom || rec.Type() == recordTypeTextBytesAtom {
			rec.RecordData = make([]byte, rec.Length())
			_, err = o.ReadAt(rec.RecordData, offset+headerSize)
			if err != nil {
				return err
			}
			if rec.Type() == recordTypeTextCharsAtom {
				textCharsAtom := TextCharsAtom(rec)
				text, err = textCharsAtom.getText()
				if err != nil {
					fmt.Printf("get text error: %v", err)
				}
			} else {
				textBytesAtom := TextBytesAtom(rec)
				text, err = textBytesAtom.getText()
				if err != nil {
					fmt.Printf("get text error: %v", err)
				}
			}

			if len(text) > 0 {
				o.texts = append(o.texts, text)
			}
		}
		offset += int64(rec.Length() + headerSize)
	}
	return nil
}
