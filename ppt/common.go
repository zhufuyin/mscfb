package ppt

import (
	"golang.org/x/text/encoding/unicode"
	"io"
)

var (
	utf16Decoder = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
)

func readRecordData(r io.ReaderAt, record *Record, offset int64) error {
	record.RecordData = make([]byte, record.DataLength)
	_, err := r.ReadAt(record.RecordData, offset+headerSize)
	if err != nil {
		return err
	}
	return nil
}
