package doc

import (
	"errors"
	"github.com/zhufuyin/mscfb/global"
)

type TextPiece struct {
	usesUnicode bool
	pd          *PieceDescriptor
	text        string
}

func NewTextPiece(start, end int, buf []byte, pd *PieceDescriptor) (*TextPiece, error) {
	tp := &TextPiece{
		pd: pd,
	}
	if pd.unicode {
		tp.usesUnicode = true
		global.Utf16Decoder.Reset()
		decBytes, err := global.Utf16Decoder.Bytes(buf)
		if err != nil {
			return nil, err
		}
		tp.text = string(decBytes)
	} else {
		if pd.decoder == nil {
			return nil, errors.New("decoder is empty")
		}
		decBytes, err := pd.decoder.Bytes(buf)
		if err != nil {
			return nil, err
		}
		tp.text = string(decBytes)
	}
	return tp, nil
}
