package doc

import (
	"github.com/zhufuyin/mscfb/global"
	"io"
)

type FibRgW97 struct {
	field_1_cbMac       int
	field_2_reserved1   int
	field_3_reserved2   int
	field_4_ccpText     int
	field_5_ccpFtn      int
	field_6_ccpHdd      int
	field_7_reserved3   int
	field_8_ccpAtn      int
	field_9_ccpEdn      int
	field_10_ccpTxbx    int
	field_11_ccpHdrTxbx int
	field_12_reserved4  int
	field_13_reserved5  int
	field_14_reserved6  int
	field_15_reserved7  int
	field_16_reserved8  int
	field_17_reserved9  int
	field_18_reserved10 int
	field_19_reserved11 int
	field_20_reserved12 int
	field_21_reserved13 int
	field_22_reserved14 int
}

func NewFibRgW97(r io.ReaderAt, offset int64) (*FibRgW97, error) {
	fibRgW97 := &FibRgW97{}
	cbMac, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgW97.field_1_cbMac = int(cbMac)
	offset += 4
	reserved1, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgW97.field_2_reserved1 = int(reserved1)
	offset += 4
	reserved2, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgW97.field_3_reserved2 = int(reserved2)
	offset += 4
	ccpText, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgW97.field_4_ccpText = int(ccpText)
	offset += 4
	ccpFtn, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgW97.field_5_ccpFtn = int(ccpFtn)
	offset += 4
	ccpHdd, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgW97.field_6_ccpHdd = int(ccpHdd)
	offset += 4
	reserved3, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgW97.field_7_reserved3 = int(reserved3)
	offset += 4
	ccpAtn, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgW97.field_8_ccpAtn = int(ccpAtn)
	offset += 4
	ccpEdn, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgW97.field_9_ccpEdn = int(ccpEdn)
	offset += 4
	ccpTxbx, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgW97.field_10_ccpTxbx = int(ccpTxbx)
	offset += 4
	ccpHdrTxbx, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgW97.field_11_ccpHdrTxbx = int(ccpHdrTxbx)
	offset += 4
	reserved4, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgW97.field_12_reserved4 = int(reserved4)
	offset += 4
	reserved5, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgW97.field_13_reserved5 = int(reserved5)
	offset += 4
	reserved6, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgW97.field_14_reserved6 = int(reserved6)
	offset += 4
	reserved7, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgW97.field_15_reserved7 = int(reserved7)
	offset += 4
	reserved8, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgW97.field_16_reserved8 = int(reserved8)
	offset += 4
	reserved9, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgW97.field_17_reserved9 = int(reserved9)
	offset += 4
	reserved10, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgW97.field_18_reserved10 = int(reserved10)
	offset += 4
	reserved11, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgW97.field_19_reserved11 = int(reserved11)
	offset += 4
	reserved12, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgW97.field_20_reserved12 = int(reserved12)
	offset += 4
	reserved13, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgW97.field_21_reserved13 = int(reserved13)
	offset += 4
	reserved14, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgW97.field_22_reserved14 = int(reserved14)

	return fibRgW97, nil
}
