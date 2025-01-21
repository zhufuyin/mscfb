package ppt

import (
	"io"
)

func readRecordData(r io.ReaderAt, record *Record, offset int64) error {
	record.RecordData = make([]byte, record.DataLength)
	_, err := r.ReadAt(record.RecordData, offset+headerSize)
	if err != nil {
		return err
	}
	return nil
}
