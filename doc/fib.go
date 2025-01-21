package doc

import (
	"fmt"
	"github.com/zhufuyin/mscfb/global"
	"io"
)

type Fib struct {
	fibBase      *FibBase
	csw          uint16
	cslw         uint16
	fibRgW97     *FibRgW97
	cbRgFcLcb97  uint16
	fibRgFcLcb97 *FibRgFcLcb97
	cswNew       uint16
	nFibNew      int16
	fibRgCswNew  []byte
}

type FibBase struct {
	field_1_wIdent     uint16
	field_2_nFib       uint16
	field_3_unused     uint16
	field_4_lid        uint16
	field_5_pnNext     uint16
	field_6_flags1     uint16
	field_7_nFibBack   uint16
	field_8_lKey       uint32
	field_9_envr       byte
	field_10_flags2    byte
	field_11_Chs       uint16
	field_12_chsTables uint16
	field_13_fcMin     uint32
	field_14_fcMac     uint32
}

func NewFib(docStream io.ReaderAt) (*Fib, error) {
	f := &Fib{}
	err := f.readFibBase(docStream)
	if err != nil {
		return nil, err
	}
	offset := int64(32)
	csw, err := global.ReadUint16At(docStream, offset)
	if err != nil {
		return nil, err
	}
	f.csw = csw
	offset += 2
	fibRw97, err := NewFibRgW97(docStream, offset)
	if err != nil {
		return nil, err
	}
	f.fibRgW97 = fibRw97
	offset = 62
	cslw, err := global.ReadUint16At(docStream, offset)
	if err != nil {
		return nil, err
	}
	f.cslw = cslw
	offset += 2
	if f.fibBase.field_2_nFib < 105 {
		// todo office 95
		return f, nil
	}
	// office 97
	fibRgLw97, err := NewFibRgW97(docStream, offset)
	if err != nil {
		return nil, err
	}
	f.fibRgW97 = fibRgLw97
	offset = 152
	cbRgFcLcb, err := global.ReadUint16At(docStream, offset)
	if err != nil {
		return nil, err
	}
	f.cbRgFcLcb97 = cbRgFcLcb
	offset += 2
	// skip fibRgFcLcbBlob
	fibRgFcLcb97, err := NewFibRgFcLcb97(docStream, offset)
	if err != nil {
		return nil, err
	}
	f.fibRgFcLcb97 = fibRgFcLcb97
	offset += int64(cbRgFcLcb * 4 * 2)
	cswNew, err := global.ReadUint16At(docStream, offset)
	if err != nil {
		return nil, err
	}
	f.cswNew = cswNew
	offset += 2
	if cswNew != 0 { // fibRgCswNew
		// version number of the file format: 0x00D9, 0x0101, 0x010C, 0x0112
		nFibNew, err := global.ReadUint16At(docStream, offset)
		if err != nil {
			return nil, err
		}
		f.nFibNew = int16(nFibNew)
		offset += 2
		// the first two bytes are stored as nFibNew
		fibRgCswNewLength := (cswNew - 1) * 2
		fibRgCswNew := make([]byte, fibRgCswNewLength)
		_, err = docStream.ReadAt(fibRgCswNew, offset)
		if err != nil {
			return nil, err
		}
		f.fibRgCswNew = fibRgCswNew
	} else { // not present in the file
		f.nFibNew = -1
		f.fibRgCswNew = []byte{}
	}
	f.assertCbRgFcLcb()
	f.assertCswNew()
	return f, nil
}

func (f *Fib) readFibBase(docStream io.ReaderAt) error {
	offset := int64(0)
	fibBase := &FibBase{}
	wIdent, err := global.ReadUint16At(docStream, offset)
	if err != nil {
		return err
	}
	fibBase.field_1_wIdent = wIdent
	offset += 2 // 0x02
	nFib, err := global.ReadUint16At(docStream, offset)
	if err != nil {
		return err
	}
	fibBase.field_2_nFib = nFib
	offset += 2 // 0x04
	unused, err := global.ReadUint16At(docStream, offset)
	if err != nil {
		return err
	}
	fibBase.field_3_unused = unused
	offset += 2 // 0x06
	lid, err := global.ReadUint16At(docStream, offset)
	if err != nil {
		return err
	}
	fibBase.field_4_lid = lid
	offset += 2 // // 0x08
	pnNext, err := global.ReadUint16At(docStream, offset)
	if err != nil {
		return err
	}
	fibBase.field_5_pnNext = pnNext
	offset += 2 // 0x0a
	flags1, err := global.ReadUint16At(docStream, offset)
	if err != nil {
		return err
	}
	fibBase.field_6_flags1 = flags1
	offset += 2 // 0x0c
	nfibBack, err := global.ReadUint16At(docStream, offset)
	if err != nil {
		return err
	}
	fibBase.field_7_nFibBack = nfibBack
	offset += 2 // 0x0e
	lKey, err := global.ReadUint32At(docStream, offset)
	if err != nil {
		return err
	}
	fibBase.field_8_lKey = lKey
	offset += 4 // 0x12
	envr, err := global.ReadByteAt(docStream, offset)
	if err != nil {
		return err
	}
	fibBase.field_9_envr = envr
	offset += 1 // 0x13
	flags2, err := global.ReadByteAt(docStream, offset)
	if err != nil {
		return err
	}
	fibBase.field_10_flags2 = flags2
	offset += 1 // 0x14
	chs, err := global.ReadUint16At(docStream, offset)
	if err != nil {
		return err
	}
	fibBase.field_11_Chs = chs
	offset += 2 // 0x16
	chsTables, err := global.ReadUint16At(docStream, offset)
	if err != nil {
		return err
	}
	fibBase.field_12_chsTables = chsTables
	offset += 2 // 0x18
	fcMin, err := global.ReadUint32At(docStream, offset)
	if err != nil {
		return err
	}
	fibBase.field_13_fcMin = fcMin
	offset += 4 // 0x1c
	fcMac, err := global.ReadUint32At(docStream, offset)
	if err != nil {
		return err
	}
	fibBase.field_14_fcMac = fcMac
	f.fibBase = fibBase
	return err
}

func (f *Fib) getNFib() int16 {
	if f.cswNew == 0 {
		return int16(f.fibBase.field_2_nFib)
	}
	return f.nFibNew
}

// MS-DOC section 2.5.1
func (f *Fib) assertCbRgFcLcb() {
	nfib := f.getNFib()
	nFibHex := fmt.Sprintf("%04X", nfib)
	switch nfib {
	case 0x00BE, 0x00BF, 0x00C0, 0x00C1, 0x00C2, 0x00C3:
		if f.cbRgFcLcb97 != 0x005D {
			fmt.Printf("Since FIB.nFib == %s value of FIB.cbRgFcLcb MUST be %s, not %x\n",
				nFibHex, "0x005D", f.cbRgFcLcb97)
		}
	case 0x00D8, 0x00D9: //  Docs "official"
		if f.cbRgFcLcb97 != 0x006C {
			fmt.Printf("Since FIB.nFib == %s value of FIB.cbRgFcLcb MUST be %s, not %x\n",
				nFibHex, "0x006C", f.cbRgFcLcb97)
		}
	case 0x0101:
		if f.cbRgFcLcb97 != 0x0088 {
			fmt.Printf("Since FIB.nFib == %s value of FIB.cbRgFcLcb MUST be %s, not %x\n",
				nFibHex, "0x0088", f.cbRgFcLcb97)
		}
		break
	case 0x010B, 0x010C: //  Docs "official"
		if f.cbRgFcLcb97 != 0x00A4 {
			fmt.Printf("Since FIB.nFib == %s value of FIB.cbRgFcLcb MUST be %s, not %x\n",
				nFibHex, "0x00A4", f.cbRgFcLcb97)
		}
		break
	case 0x0112:
		if f.cbRgFcLcb97 != 0x00B7 {
			fmt.Printf("Since FIB.nFib == %s value of FIB.cbRgFcLcb MUST be %s, not %x\n",
				nFibHex, "0x00B7", f.cbRgFcLcb97)
		}
		break
	default:
		/* The Word spec has a much smaller list of "valid" values
		 * to what the large CommonCrawl corpus contains!
		 */
		fmt.Printf("Invalid file format version number: %d(%s)\n", nfib, nFibHex)
	}
}

// MS-DOC section 2.5.1
func (f *Fib) assertCswNew() {
	nfib := f.getNFib()
	nfibHex := fmt.Sprintf("%04X", nfib)
	cswNewHex := fmt.Sprintf("%04X", f.cswNew)
	switch nfib {
	case 0x00C1:
		if f.cswNew != 0x0000 {
			fmt.Printf("Since FIB.nFib == %s value of FIB.cswNew MUST be %s, not %s\n",
				nfibHex, "0x0000", cswNewHex)
		}
	case 0x00D9, 0x0101, 0x010C:
		if f.cswNew != 0x0002 {
			fmt.Printf("Since FIB.nFib == %s value of FIB.cswNew MUST be %s, not %s\n",
				nfibHex, "0x0002", cswNewHex)
		}
	case 0x0112:
		if f.cswNew != 0x0005 {
			fmt.Printf("Since FIB.nFib == %s value of FIB.cswNew MUST be %s, not %s\n",
				nfibHex, "0x0005", cswNewHex)
		}
	default:
		fmt.Printf("Invalid file format version number: %s\n", nfibHex)
	}
}

func (fb *FibBase) isFWhichTblStm() bool {
	mask := uint16(0x0200)
	fWhichTblStm := mask & fb.field_6_flags1
	if fWhichTblStm != 0 {
		return true
	}
	return false
}
