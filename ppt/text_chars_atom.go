package ppt

import "github.com/zhufuyin/mscfb/global"

type TextCharsAtom Record

func (a *TextCharsAtom) getText() (string, error) {
	global.Utf16Decoder.Reset()
	decBytes, err := global.Utf16Decoder.Bytes(a.RecordData)
	if err != nil {
		return "", err
	}
	return string(decBytes), nil
}
