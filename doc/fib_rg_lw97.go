package doc

import (
	"github.com/zhufuyin/mscfb/global"
	"io"
)

type FibRgLw97 struct {
	field_1_cbMac       uint32
	field_2_reserved1   uint32
	field_3_reserved2   uint32
	field_4_ccpText     uint32
	field_5_ccpFtn      uint32
	field_6_ccpHdd      uint32
	field_7_reserved3   uint32
	field_8_ccpAtn      uint32
	field_9_ccpEdn      uint32
	field_10_ccpTxbx    uint32
	field_11_ccpHdrTxbx uint32
	field_12_reserved4  uint32
	field_13_reserved5  uint32
	field_14_reserved6  uint32
	field_15_reserved7  uint32
	field_16_reserved8  uint32
	field_17_reserved9  uint32
	field_18_reserved10 uint32
	field_19_reserved11 uint32
	field_20_reserved12 uint32
	field_21_reserved13 uint32
	field_22_reserved14 uint32
}

func NewFibRgLw97(r io.ReaderAt, offset int64) (*FibRgLw97, error) {
	fibRgLw97 := &FibRgLw97{}
	cbMac, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgLw97.field_1_cbMac = cbMac
	offset += 4 // 0x04
	reserved1, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgLw97.field_2_reserved1 = reserved1
	offset += 4
	reserved2, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgLw97.field_3_reserved2 = reserved2
	offset += 4
	cppText, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgLw97.field_4_ccpText = cppText
	offset += 4
	ccpFtn, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgLw97.field_5_ccpFtn = ccpFtn
	offset += 4
	cppHdd, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgLw97.field_6_ccpHdd = cppHdd
	offset += 4
	reserved3, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgLw97.field_7_reserved3 = reserved3
	offset += 4
	cppAtn, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgLw97.field_8_ccpAtn = cppAtn
	offset += 4
	cppEdn, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgLw97.field_9_ccpEdn = cppEdn
	offset += 4
	cppTxbx, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgLw97.field_10_ccpTxbx = cppTxbx
	offset += 4
	cppHdrTxbx, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgLw97.field_11_ccpHdrTxbx = cppHdrTxbx
	offset += 4
	reserved4, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgLw97.field_12_reserved4 = reserved4
	offset += 4
	reserved5, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgLw97.field_13_reserved5 = reserved5
	offset += 4
	reserved6, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgLw97.field_14_reserved6 = reserved6
	offset += 4
	reserved7, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgLw97.field_15_reserved7 = reserved7
	offset += 4
	reserved8, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgLw97.field_16_reserved8 = reserved8
	offset += 4
	reserved9, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgLw97.field_17_reserved9 = reserved9
	offset += 4
	reserved10, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgLw97.field_18_reserved10 = reserved10
	offset += 4
	reserved11, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgLw97.field_19_reserved11 = reserved11
	offset += 4
	reserved12, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgLw97.field_20_reserved12 = reserved12
	offset += 4
	reserved13, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgLw97.field_21_reserved13 = reserved13
	offset += 4
	reserved14, err := global.ReadUint32At(r, offset)
	if err != nil {
		return nil, err
	}
	fibRgLw97.field_22_reserved14 = reserved14

	return fibRgLw97, nil

}
