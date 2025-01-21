package doc

import (
	"io"
)

type TextPieceTable struct {
	cpMin      int32
	textPieces []*TextPiece
}

func NewTextPieceTable(docStream, tblStream io.ReaderAt, offset, size int) (*TextPieceTable, error) {
	tpt := &TextPieceTable{}
	pieceTable, err := NewPlexOfCps(tblStream, offset, size, 8)
	if err != nil {
		return nil, err
	}
	pds := []*PieceDescriptor{}
	for i := 0; i < pieceTable.iMac; i++ {
		prop := pieceTable.props[i]
		// parse pcd
		pd, err := NewPieceDescriptor(prop.buf, 0)
		if err != nil {
			return nil, err
		}
		pds = append(pds, pd)
	}
	tpt.cpMin = pds[0].fc
	for _, piece := range pds {
		if piece.fc < tpt.cpMin {
			tpt.cpMin = piece.fc
		}
	}
	// build list of TextPieces
	for i, pd := range pds {
		start := pd.fc
		prop := pieceTable.props[i]
		cpStart, cpEnd := prop.cpStart, prop.cpEnd
		multiple := 1
		if pd.unicode {
			multiple = 2
		}
		textSizeBytes := (cpEnd - cpStart) * multiple
		buf := make([]byte, textSizeBytes)
		_, err = docStream.ReadAt(buf, int64(start))
		if err != nil {
			return nil, err
		}
		tp, err := NewTextPiece(cpStart, cpEnd, buf, pd)
		if err != nil {
			return nil, err
		}
		tpt.textPieces = append(tpt.textPieces, tp)
	}
	return tpt, nil
}
