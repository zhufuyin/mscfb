package doc

import (
	"fmt"
	"github.com/zhufuyin/mscfb/global"
	"io"
)

type ComplexFileTable struct {
	tpt *TextPieceTable
}

// parse Clx with offset from FibRgFcLcb97.fcClx
func NewComplexFileTable(docStream, tableStream io.ReaderAt, fcClx int64) (*ComplexFileTable, error) {
	offset := fcClx
	// RgPrc array part
	for {
		clxtBytes := make([]byte, 1)
		_, err := tableStream.ReadAt(clxtBytes, offset)
		//clxt, err := global.ReadByteAt(tableStream, offset)
		if err != nil {
			return nil, err
		}
		clxt := clxtBytes[0]
		if clxt != 0x01 { // Pcd part tag
			break
		}
		offset++
		size, err := global.ReadUint16At(tableStream, offset) // size
		if err != nil {
			return nil, err
		}
		offset += int64(2 + size)
	}
	// Pcdt part
	clxt, err := global.ReadByteAt(tableStream, offset)
	if err != nil {
		return nil, err
	}
	if clxt != 0x02 {
		fmt.Printf("The text piece table is corrupted, expected byte value 0x02 but had %x\n", clxt)
		return nil, fmt.Errorf("The text piece table is corrupted, expected byte value 0x02 but had %x", clxt)
	}
	offset++
	pieceTableSize, err := global.ReadUint32At(tableStream, offset)
	if err != nil {
		return nil, err
	}
	offset += 4
	// handle PlcPcd, that is text piece table
	tpt, err := NewTextPieceTable(docStream, tableStream, int(offset), int(pieceTableSize))
	if err != nil {
		return nil, err
	}
	return &ComplexFileTable{tpt: tpt}, nil
}
