package ppt

type CurrentUserAtom struct {
	Record
	size                uint32
	headerToken         uint32
	offsetToCurrentEdit uint32
	lenUserName         uint16
	docFileVersion      uint16
	majorVersion        uint8
	minorVersion        uint8
	unused              uint16
	ansiUserName        PrintableAnsiString
	relVersion          uint32
	unicodeUserName     PrintableUnicodeString
}

type PrintableAnsiString struct {
}

type PrintableUnicodeString struct{}

func (a *CurrentUserAtom) parse() {
	offset := 0
	a.size = a.Uint32At(offset)
	offset += 4
	a.headerToken = a.Uint32At(offset)
	offset += 4
	a.offsetToCurrentEdit = a.Uint32At(offset)
	offset += 4
	a.lenUserName = a.Uint16At(offset)
	offset += 2
	a.docFileVersion = a.Uint16At(offset)
	offset += 2
	a.majorVersion = a.ByteAt(offset)
	offset += 1
	a.minorVersion = a.ByteAt(offset)
}
