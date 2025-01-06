package ppt

type OfficeArtSpContainer struct {
	Record
}

func (c *OfficeArtSpContainer) readClientTextbox() (*OfficeArtClientTextbox, error) {
	offset := int64(0)
	for offset < c.DataLength {
		record, err := readRecordHeaderOnly(c, offset, recordTypeUnspecified)
		if err != nil {
			return nil, err
		}
		if record.RecType == recordTypeOfficeArtClientTextbox {
			record.RecordData = make([]byte, record.DataLength)
			_, err := c.ReadAt(record.RecordData, offset+headerSize)
			if err != nil {
				return nil, err
			}
			textbox := OfficeArtClientTextbox{
				Record: record,
			}
			return &textbox, nil
		}
		offset += record.DataLength + headerSize
	}
	return nil, nil
}
