package ppt

type PersistDirectoryAtom struct {
	Record
	rgPersistDirEntry []*PersistDirectoryEntry
}

func (a *PersistDirectoryAtom) parse() error {
	const persistOffsetEntrySize = 4
	offset := 0
	for offset < int(a.DataLength) {
		persistDirectoryEntry, err := parsePersistDirectoryEntry(a.RecordData, offset)
		if err != nil {
			return err
		}
		a.rgPersistDirEntry = append(a.rgPersistDirEntry, persistDirectoryEntry)
		offset += 4 + int(persistDirectoryEntry.cPersist*persistOffsetEntrySize)
	}
	return nil
}
