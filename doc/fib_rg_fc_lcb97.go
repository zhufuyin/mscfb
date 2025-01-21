package doc

import (
	"encoding/binary"
	"io"
)

const (
	fcClxOffset = 66 * 4
)

type FibRgFcLcb97 struct {
	fcClx uint32 // specify an offset in the table stream to the Clx beginning
}

func NewFibRgFcLcb97(docStream io.ReaderAt, offset int64) (*FibRgFcLcb97, error) {
	fibRgFcLcb97 := &FibRgFcLcb97{}
	fcClxBytes := make([]byte, 4)
	_, err := docStream.ReadAt(fcClxBytes, offset+fcClxOffset)
	if err != nil {
		return nil, err
	}
	fibRgFcLcb97.fcClx = binary.LittleEndian.Uint32(fcClxBytes)
	return fibRgFcLcb97, nil
}
