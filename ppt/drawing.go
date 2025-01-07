package ppt

import (
	"strings"
)

// read text from DrawingContainer, refer [MS-ODRAW] 2.2.13
func readTextFromDrawing(drawing Record) (string, error) {
	//fmt.Printf("drawing type=%X\n", drawing.Type())
	// read OfficeArtDgContainer from drawing data
	officeArtDgContainerRecord, err := readRecord(drawing, 0, recordTypeOfficeArtDgContainer)
	if err != nil {
		return "", err
	}
	// read shape part (OfficeArtSpContainer) in OfficeArtDgContainer
	officeArtDgContainer := &OfficeArtDGContainer{
		Record: officeArtDgContainerRecord,
	}
	officeArtDgContainer.parse()
	if len(officeArtDgContainer.spContainerRecords) == 0 {
		return "", nil
	}

	// read clientTextbox part ([MS-PPT 2.9.76] OfficeArtClientTextbox) in OfficeArtSpContainer
	var texts []string
	for _, sp := range officeArtDgContainer.spContainerRecords {
		err := sp.parse()
		if err != nil {
			continue
		}
		txts := sp.extractText()
		if len(txts) > 0 {
			texts = append(texts, txts...)
		}
	}
	return strings.Join(texts, "\n"), nil

	//officeArtSpContainerRecord, err := readTargetRecordInstance(officeArtDgContainerRecord, 0,
	//	InstanceType{
	//		Instance: 0x000,
	//		Type:     0xF004,
	//	})
	//if err != nil {
	//	return "", err
	//}

	//officeArtClientTextboxRecord, err := readTargetRecordInstance(officeArtSpContainerRecord, 0,
	//	InstanceType{
	//		Instance: 0x000,
	//		Type:     0xF00D,
	//	})
	//if err != nil {
	//	return "", err
	//}
	//// read TextClientDataSubContainerOrAtom record array and then filter RT_TextCharsAtom and RT_TextBytesAtom
	//offset := int64(headerSize)
	//textBuilder := &strings.Builder{}
	//for {
	//	rec, err := readRecordHeaderOnly(officeArtClientTextboxRecord, offset, recordTypeUnspecified)
	//	if err != nil {
	//		return "", err
	//	}
	//	fmt.Printf("record type=%X instance=%X\n", rec.Type(), rec.Instance())
	//	if rec.Type() == recordTypeTextCharsAtom {
	//		rec.RecordData = make([]byte, rec.Length())
	//		_, err = officeArtClientTextboxRecord.ReadAt(rec.RecordData, offset+headerSize)
	//		if err != nil {
	//			return "", err
	//		}
	//		err = readTextFromTextCharsAtom(rec, textBuilder, utf16Decoder)
	//		if err != nil {
	//			return "", err
	//		}
	//	} else if rec.Type() == recordTypeTextBytesAtom {
	//		rec.RecordData = make([]byte, rec.Length())
	//		_, err = officeArtClientTextboxRecord.ReadAt(rec.RecordData, offset+headerSize)
	//		if err != nil {
	//			return "", err
	//		}
	//		err = readTextFromTextBytesAtom(rec, textBuilder, utf16Decoder)
	//		if err != nil {
	//			return "", err
	//		}
	//	}
	//	offset += int64(rec.Length() + headerSize)
	//	if offset >= int64(officeArtClientTextboxRecord.Length()) {
	//		break
	//	}
	//}
	//return textBuilder.String(), nil
}
