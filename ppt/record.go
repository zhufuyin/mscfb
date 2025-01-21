package ppt

import (
	"encoding/binary"
	"errors"
	"github.com/zhufuyin/mscfb/global"
	"io"
)

const headerSize = 8

// recordType is an enumeration that specifies the record type of an atom record or a container record
// ([MS-PPT] 2.13.24 RecordType)
type RecordType uint16

const (
	recordTypeUnspecified              RecordType = 0
	recordTypeDocument                 RecordType = 0x03E8
	recordTypeSlide                    RecordType = 0x03EE
	recordTypeNotes                    RecordType = 0x03F0
	recordTypeEnvironment              RecordType = 0x03F2
	recordTypeSlidePersistAtom         RecordType = 0x03F3
	recordTypeMainMaster               RecordType = 0x03F8
	recordTypeSlideShowSlideInfoAtom   RecordType = 0x03F9
	recordTypeExternalObjectList       RecordType = 0x0409
	recordTypeDrawingGroup             RecordType = 0x040B
	recordTypeDrawing                  RecordType = 0x040C
	recordTypeList                     RecordType = 0x07D0
	recordTypeSoundCollection          RecordType = 0x07E4
	recordTypeTextCharsAtom            RecordType = 0x0FA0
	recordTypeTextBytesAtom            RecordType = 0x0FA8
	recordTypeHeadersFooters           RecordType = 0x0FD9
	recordTypeSlideListWithText        RecordType = 0x0FF0
	recordTypeUserEditAtom             RecordType = 0x0FF5
	recordTypeCurrentUserAtom          RecordType = 0x0FF6
	recordTypePersistDirectoryAtom     RecordType = 0x1772
	recordTypeRoundTripSlideSyncInfo12 RecordType = 0x3714
	// recordType within [MS-ODRAW]
	recordTypeOfficeArtDgContainer   RecordType = 0xF002
	recordTypeOfficeArtSpgrContainer RecordType = 0xF003
	recordTypeOfficeArtSpContainer   RecordType = 0xF004
	recordTypeOfficeArgFdg           RecordType = 0xF008
	recordTypeOfficeArtClientTextbox RecordType = 0xF00D
)

// LowerPart returns lower byte of record type
func (r RecordType) LowerPart() byte {
	const fullByte = 0xFF
	return byte(r & fullByte)
}

var errMismatchRecordType = errors.New("mismatch record type")

type Record struct {
	header [headerSize]byte
	RecordData
	offset       int64
	RecType      RecordType
	TypeInstance uint16
	DataLength   int64
}

// Version of record data
func (r Record) Version() uint16 {
	data := binary.LittleEndian.Uint16(r.header[:2])
	versionMask := uint16(0xF)
	version := data & versionMask
	return version
}

// Instance specifies the record instance data
func (r Record) Instance() uint16 {
	data := binary.LittleEndian.Uint16(r.header[:2])
	//instanceMask := uint16(0xFFF0)
	instance := data >> 4
	return instance
}

// Type returns recordType of record contained in it's header
func (r Record) Type() RecordType {
	return RecordType(binary.LittleEndian.Uint16(r.header[2:4]))
}

// Length returns data length contained in record header
func (r Record) Length() uint32 {
	return binary.LittleEndian.Uint32(r.header[4:8])
}

// Data returns all data from record except header
func (r Record) Data() []byte {
	return r.RecordData
}

type RecordData []byte

// ReadAt copies bytes from record data at given offset into buffer p
func (rd RecordData) ReadAt(p []byte, off int64) (n int, err error) {
	return global.ReadBytes(rd, p, off)
}

func (rd RecordData) Uint64At(offset int) uint64 {
	return binary.LittleEndian.Uint64(rd[offset:])
}

// Uint32At interprets 4 bytes of record data at given offset as uint32 value and returns it
func (rd RecordData) Uint32At(offset int) uint32 {
	return binary.LittleEndian.Uint32(rd[offset:])
}

func (rd RecordData) Uint16At(offset int) uint16 {
	return binary.LittleEndian.Uint16(rd[offset:])
}

func (rd RecordData) ByteAt(offset int) uint8 {
	return rd[offset]
}

// readRecord reads header and data of record. If wantedType is specified (not equals recordTypeUnspecified),
// also compares read type with the wanted one and returns an error is they are not equal
func readRecord(f io.ReaderAt, offset int64, wantedType RecordType) (Record, error) {
	r, err := readRecordHeaderOnly(f, offset, wantedType)
	if err != nil {
		return Record{}, err
	}
	r.RecordData = make([]byte, r.Length())
	_, err = f.ReadAt(r.RecordData, offset+headerSize)
	if err != nil {
		return Record{}, err
	}
	return r, nil
}

// readRecordHeaderOnly reads header of record. If wantedType is specified (not equals recordTypeUnspecified),
// also compares read type with the wanted one and returns an error is they are not equal
func readRecordHeaderOnly(f io.ReaderAt, offset int64, wantedType RecordType) (Record, error) {
	r := Record{}
	_, err := f.ReadAt(r.header[:], offset)
	if err != nil {
		return Record{}, err
	}
	r.offset = offset
	r.RecType = r.Type()
	r.TypeInstance = r.Instance()
	r.DataLength = int64(r.Length())
	if wantedType != recordTypeUnspecified && r.RecType != wantedType {
		//fmt.Printf("recordType=%X  wantedType=%X\n", r.Type(), wantedType)
		return Record{}, errMismatchRecordType
	}
	return r, nil
}
