package ppt

import "context"

// [MS-PPT] 2.4.14.1
type MasterListWithTextContainer struct {
	Record
	rgMasterPersistAtom []*MasterPersistAtom
}

func (c *MasterListWithTextContainer) parse(ctx context.Context) error {
	offset := int64(0)
	for offset < c.DataLength {
		record, err := readRecord(c, offset, recordTypeSlidePersistAtom)
		if err != nil {
			return err
		}
		masterPersistAtom := &MasterPersistAtom{
			Record: record,
		}
		err = masterPersistAtom.parse(ctx)
		if err != nil {
			return err
		}
		c.rgMasterPersistAtom = append(c.rgMasterPersistAtom, masterPersistAtom)
		offset += record.DataLength + headerSize
	}
	return nil
}
