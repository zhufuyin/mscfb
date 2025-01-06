package ppt

type PersistDirectoryEntry struct {
	persistId       uint32
	cPersist        uint32
	rgPersistOffset []uint32
	persistIdOffset map[uint32]int64
}

// read PersistDirectoryEntry within PersistDirectoryAtom
func parsePersistDirectoryEntry(data RecordData, offset int) (*PersistDirectoryEntry, error) {
	persistDirectoryEntry := &PersistDirectoryEntry{
		persistIdOffset: make(map[uint32]int64),
	}
	persist := data.Uint32At(offset)
	persistId := persist & 0x000FFFFF
	cPersist := persist >> 20
	persistDirectoryEntry.persistId = persistId
	persistDirectoryEntry.cPersist = cPersist
	offset += 4
	for i := uint32(0); i < cPersist; i++ {
		persistOffset := data.Uint32At(offset)
		persistDirectoryEntry.rgPersistOffset = append(persistDirectoryEntry.rgPersistOffset, persistOffset)
		persistDirectoryEntry.persistIdOffset[persistId+i] = int64(persistOffset)
		offset += 4
	}
	return persistDirectoryEntry, nil
}
