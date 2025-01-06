package ppt

type TextBytesAtom Record

func (a *TextBytesAtom) getText() (string, error) {
	utf16Decoder.Reset()
	decBytes, err := decodeTextBytesAtom(a.RecordData, utf16Decoder)
	if err != nil {
		return "", err
	}
	return string(decBytes), nil
}
