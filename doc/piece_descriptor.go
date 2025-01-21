package doc

import (
	"encoding/binary"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

type PieceDescriptor struct {
	descriptor uint16
	fc         int32
	unicode    bool
	decoder    *encoding.Decoder
}

func NewPieceDescriptor(buf []byte, offset int64) (*PieceDescriptor, error) {
	pd := &PieceDescriptor{}
	descriptor := binary.LittleEndian.Uint16(buf[offset : offset+2])
	pd.descriptor = descriptor
	offset += 2
	fc := binary.LittleEndian.Uint32(buf[offset : offset+4])
	pd.fc = int32(fc)
	offset += 4
	if fc&0x40000000 == 0 {
		pd.unicode = true
	} else {
		pd.unicode = false
		pd.fc = int32(fc) & ^int32(0x40000000)
		pd.fc /= 2
		pd.decoder = charmap.Windows1252.NewDecoder()
	}
	return pd, nil
}
