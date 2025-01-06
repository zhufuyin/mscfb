package ppt

type UserEditAtom struct {
	Record
	lastSlideIdRef         uint32
	version                uint16
	minorVersion           uint8
	majorVersion           uint8
	offsetLastEdit         uint32
	offsetPersistDirectory uint32
	docPersistIdRef        uint32
	persistIdSeed          uint32
	lastView               uint16
}

func (a *UserEditAtom) parse() {
	offset := 0
	a.lastSlideIdRef = a.Uint32At(offset)
	offset += 4
	a.version = a.Uint16At(offset)
	offset += 2
	a.minorVersion = a.ByteAt(offset)
	offset += 1
	a.majorVersion = a.ByteAt(offset)
	offset += 1
	a.offsetLastEdit = a.Uint32At(offset)
	offset += 4
	a.offsetPersistDirectory = a.Uint32At(offset)
}
