package ppt

type OfficeArtSpgrContainer struct {
	Record
	spContainerRecords []OfficeArtSpContainer
}

func (g *OfficeArtSpgrContainer) readOfficeArtSpContainer() error {
	offset := int64(0)
	for offset < g.DataLength {
		record, err := readRecord(g, offset, recordTypeUnspecified)
		if err != nil {
			return err
		}
		switch record.RecType {
		case recordTypeOfficeArtSpContainer:
			g.spContainerRecords = append(g.spContainerRecords, OfficeArtSpContainer{record})
		case recordTypeOfficeArtSpgrContainer:
			subGroup := &OfficeArtSpgrContainer{
				Record: record,
			}
			subGroup.readOfficeArtSpContainer()
			if len(subGroup.spContainerRecords) > 0 {
				g.spContainerRecords = append(g.spContainerRecords, subGroup.spContainerRecords...)
			}
		}
		offset += int64(record.Length() + headerSize) // header + data length
	}
	return nil
}
