package ppt

// [MS-ODRAW] 2.2.14
type OfficeArtSpContainer struct {
	Record
	clientTextbox *OfficeArtClientTextbox
}

func (c *OfficeArtSpContainer) parse() error {
	offset := int64(0)
	for offset < c.DataLength {
		record, err := readRecordHeaderOnly(c, offset, recordTypeUnspecified)
		if err != nil {
			return err
		}
		if record.RecType == recordTypeOfficeArtClientTextbox {
			record.RecordData = make([]byte, record.DataLength)
			_, err := c.ReadAt(record.RecordData, offset+headerSize)
			if err != nil {
				return err
			}
			textbox := &OfficeArtClientTextbox{
				Record: record,
			}
			err = textbox.parse()
			if err != nil {
				return err
			}
			c.clientTextbox = textbox
			return nil
		}
		offset += record.DataLength + headerSize
	}
	return nil
}

func (c *OfficeArtSpContainer) extractText() []string {
	if c.clientTextbox == nil {
		return nil
	}
	return c.clientTextbox.texts
}
