package ppt

type OfficeArtDGContainer struct {
	Record
	spContainerRecords []OfficeArtSpContainer
}

func (c *OfficeArtDGContainer) readOfficeArtSpContainers() error {
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
			c.spContainerRecords = append(c.spContainerRecords, OfficeArtSpContainer{record})
		case recordTypeOfficeArtSpgrContainer:
			record.RecordData = make([]byte, record.Length())
			_, err := c.ReadAt(record.RecordData, offset+headerSize)
			if err != nil {
				return err
			}
			group := &OfficeArtSpgrContainer{
				Record: record,
			}
			group.readOfficeArtSpContainer()
			if len(group.spContainerRecords) > 0 {
				c.spContainerRecords = append(c.spContainerRecords, group.spContainerRecords...)
			}
		}
		offset += int64(headerSize + record.Length())
	}
	return nil
}
