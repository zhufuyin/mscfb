package ppt

type TextCharsAtom Record

func (a TextCharsAtom) getText() (string, error) {
	utf16Decoder.Reset()
	decBytes, err := utf16Decoder.Bytes(a.RecordData)
	if err != nil {
		return "", err
	}
	return string(decBytes), nil
}
