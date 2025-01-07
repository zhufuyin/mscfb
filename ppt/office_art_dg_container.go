package ppt

type OfficeArtDGContainer struct {
	Record
	spContainerRecords []*OfficeArtSpContainer
}

func (c *OfficeArtDGContainer) parse() error {
	offset := int64(0)
	for offset < c.DataLength {
		record, err := readRecordHeaderOnly(c, offset, recordTypeUnspecified)
		if err != nil {
			return err
		}
		switch record.RecType {
		case recordTypeOfficeArtSpContainer:
			record.RecordData = make([]byte, record.Length())
			_, err := c.ReadAt(record.RecordData, offset+headerSize)
			if err != nil {
				return err
			}
			spContainer := &OfficeArtSpContainer{
				Record: record,
			}
			err = spContainer.parse()
			if err != nil {
				return err
			}
			c.spContainerRecords = append(c.spContainerRecords, spContainer)
		case recordTypeOfficeArtSpgrContainer:
			record.RecordData = make([]byte, record.Length())
			_, err := c.ReadAt(record.RecordData, offset+headerSize)
			if err != nil {
				return err
			}
			group := &OfficeArtSpgrContainer{
				Record: record,
			}
			err = group.parse()
			if err != nil {
				return err
			}
			if len(group.spContainerRecords) > 0 {
				c.spContainerRecords = append(c.spContainerRecords, group.spContainerRecords...)
			}
		}
		offset += int64(headerSize + record.Length())
	}
	for _, spContainer := range c.spContainerRecords {
		err := spContainer.parse()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *OfficeArtDGContainer) extractText() []string {
	if len(c.spContainerRecords) == 0 {
		return nil
	}
	var texts []string
	for _, spContainer := range c.spContainerRecords {
		texts = append(texts, spContainer.extractText()...)
	}
	return texts
}
