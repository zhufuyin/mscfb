package ppt

import "context"

// [MS-PPT] 2.4.14.2
type MasterPersistAtom struct {
	Record
	persistIdRef uint32 // refer to MainMasterContainer or SlideContainer
}

func (m *MasterPersistAtom) parse(ctx context.Context) error {
	m.persistIdRef = m.Uint32At(0)
	return nil
}
