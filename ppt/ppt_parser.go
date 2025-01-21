package ppt

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/zhufuyin/mscfb/cfb"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"io"
	"strings"
)

// skipped metadata or non-readable records in slide container
var slideSkippedRecordsTypes = []RecordType{
	recordTypeExternalObjectList,
	recordTypeEnvironment,
	recordTypeSoundCollection,
	recordTypeDrawingGroup,
	recordTypeSlideListWithText,
	recordTypeList,
	recordTypeHeadersFooters,
	recordTypeHeadersFooters,
}

// skipped metadata or non-readable records in drawing container
var drawingSkippedRecordsTypes = []RecordType{
	recordTypeSlideShowSlideInfoAtom,
	recordTypeHeadersFooters,
	recordTypeRoundTripSlideSyncInfo12,
}

const (
	userPersistIDRefOffset = 16
)

// ExtractText parses PPT file represented by Reader r and extracts text from it.
func ExtractText(r io.Reader) (string, error) {
	ra := global.NewReaderAt(r)

	d, err := cfb.New(ra)
	if err != nil {
		return "", err
	}
	currentUser, pptDocument := getCurrentUserAndPPTDoc(d)
	if err := isValidPPT(currentUser, pptDocument); err != nil {
		return "", err
	}
	offsetPersistDirectory, liveRecord, err := getUserEditAtomsData(currentUser, pptDocument)
	if err != nil {
		return "", err
	}
	persistDirEntries, err := getPersistDirectoryEntries(pptDocument, offsetPersistDirectory)
	if err != nil {
		return "", err
	}

	// get DocumentContainer reference
	docPersistIDRef := liveRecord.Uint32At(userPersistIDRefOffset)
	documentContainer, err := readRecord(pptDocument, persistDirEntries[docPersistIDRef], recordTypeDocument)
	if err != nil {
		return "", err
	}

	return readSlides(documentContainer, pptDocument, persistDirEntries)
}

// toMemoryBuffer transforms io.Reader to in-memory io.ReaderAt

// getCurrentUserAndPPTDoc extracts necessary cfb files from PPT file
func getCurrentUserAndPPTDoc(r *cfb.Reader) (currentUser *cfb.File, pptDocument *cfb.File) {
	for _, f := range r.File {
		switch f.Name {
		case "Current User":
			currentUser = f
		case "PowerPoint Document":
			pptDocument = f
		}
	}
	return currentUser, pptDocument
}

// isValidPPT checks if provided file is valid, meaning
// it has both "Current User" and "PowerPoint Document" files
// and "Current User"'s CurrentUserAtom record has valid header token
func isValidPPT(currentUser, pptDocument *cfb.File) error {
	const (
		headerTokenOffset      = 12
		encryptedDocumentToken = 0xF3D1C4DF
		plainDocumentToken     = 0xE391C05F
	)

	if currentUser == nil || pptDocument == nil {
		return fmt.Errorf(".ppt file must contain \"Current User\" and \"PowerPoint Document\" streams")
	}
	var b [4]byte
	_, err := currentUser.ReadAt(b[:], headerTokenOffset)
	if err != nil {
		return err
	}
	headerToken := binary.LittleEndian.Uint32(b[:])
	if headerToken != plainDocumentToken && headerToken != encryptedDocumentToken {
		return fmt.Errorf("invalid UserEditAtom header token %X", headerToken)
	}
	return nil
}

// getUserEditAtomsData extracts "live record" and persist directory offsets
// according to section 2.1.2 of specification (https://msopenspecs.azureedge.net/files/MS-PPT/%5bMS-PPT%5d-210422.pdf)
func getUserEditAtomsData(currentUser, pptDocument *cfb.File) (persistDirectoryOffsets []int64,
	liveRecord Record, err error) {
	const (
		offsetLastEditInitialPosition  = 16
		offsetLastEditPosition         = 8
		persistDirectoryOffsetPosition = 12
	)
	var b [4]byte
	_, err = currentUser.ReadAt(b[:], offsetLastEditInitialPosition)
	if err != nil {
		return nil, Record{}, err
	}
	offsetLastEdit := binary.LittleEndian.Uint32(b[:]) // 缺少判断offsetLastEdit=0的终止条件，0表示没有user edit

	for {
		liveRecord, err = readRecord(pptDocument, int64(offsetLastEdit), recordTypeUserEditAtom)
		if err != nil {
			if errors.Is(err, errMismatchRecordType) {
				break
			}
			return nil, Record{}, err
		}
		persistDirectoryOffsets = append(
			persistDirectoryOffsets,
			int64(liveRecord.Uint32At(persistDirectoryOffsetPosition)),
		) // 去掉header之后的偏移量
		offsetLastEdit = liveRecord.Uint32At(offsetLastEditPosition)
		if offsetLastEdit == 0 {
			break
		}
	}

	return persistDirectoryOffsets, liveRecord, err
}

// getPersistDirectoryEntries transforms offsets into persists directory identifiers and persist offsets according
// to section 2.1.2 of specification (https://msopenspecs.azureedge.net/files/MS-PPT/%5bMS-PPT%5d-210422.pdf)
func getPersistDirectoryEntries(pptDocument *cfb.File, offsets []int64) (map[uint32]int64, error) {
	const persistOffsetEntrySize = 4

	persistDirEntries := make(map[uint32]int64)
	for i := len(offsets) - 1; i >= 0; i-- {
		rgPersistDirEntry, err := readRecord(pptDocument, offsets[i], recordTypePersistDirectoryAtom)
		if err != nil {
			return nil, err
		}

		rgPersistDirEntryData := rgPersistDirEntry.RecordData

		for j := 0; j < len(rgPersistDirEntryData); {
			persist := rgPersistDirEntryData.Uint32At(j)
			persistID := persist & 0x000FFFFF
			cPersist := ((persist & 0xFFF00000) >> 20) & 0x00000FFF
			j += 4

			for k := uint32(0); k < cPersist; k++ {
				persistDirEntries[persistID+k] = int64(rgPersistDirEntryData.Uint32At(j + int(k)*persistOffsetEntrySize))
			}
			j += int(cPersist * persistOffsetEntrySize)
		}
	}
	return persistDirEntries, nil
}

// readSlides reads text from slides of given DocumentContainer [MS-PPT] 2.4.1
func readSlides(documentContainer Record, pptDocument io.ReaderAt, persistDirEntries map[uint32]int64) (string, error) {
	const slideSkipInitialOffset = 48 // documentAtom(48 bytes) in DocumentContainer
	//offset, err := skipRecords(documentContainer, slideSkipInitialOffset, slideSkippedRecordsTypes)
	//if err != nil {
	//	return "", err
	//}
	//slideList, err := readRecord(documentContainer, offset, recordTypeSlideListWithText)
	//if err != nil {
	//	return "", err
	//}

	slideList, err := readTargetRecordInstance(documentContainer, slideSkipInitialOffset,
		SlideListWithTextContainerType)
	if err != nil {
		return "", err
	}

	var out strings.Builder
	n := len(slideList.Data())
	for i := 0; i < n; {
		block, err := readRecord(slideList, int64(i), recordTypeUnspecified)
		if err != nil {
			return "", err
		}
		switch block.RecType {
		case recordTypeSlidePersistAtom:
			err = readTextFromSlidePersistAtom(block, pptDocument, persistDirEntries, &out)
		case recordTypeTextCharsAtom:
			err = readTextFromTextCharsAtom(block, &out, utf16Decoder)
		case recordTypeTextBytesAtom:
			err = readTextFromTextBytesAtom(block, &out, utf16Decoder)
		}
		if err != nil {
			return "", err
		}

		i += len(block.Data()) + 8
	}

	return out.String(), nil
}

func readTextFromSlidePersistAtom(
	block Record,
	pptDocument io.ReaderAt,
	persistDirEntries map[uint32]int64,
	out *strings.Builder,
) error {
	const (
		slidePersistAtomSkipInitialOffset = 32
		headerRecordTypeOffset            = 2
	)

	persistDirID := block.Uint32At(0)
	// extract slide from persist directory
	slide, err := readRecord(pptDocument, persistDirEntries[persistDirID], recordTypeSlide)
	if err != nil {
		return err
	}
	// skip metadata in SlideContainer
	offset, err := skipRecords(slide, slidePersistAtomSkipInitialOffset, drawingSkippedRecordsTypes)
	if err != nil {
		return err
	}

	drawing, err := readRecord(slide, offset, recordTypeDrawing)
	if err != nil {
		return err
	}
	text, err := readTextFromDrawing(drawing)
	if err != nil {
		return err
	}
	out.Write([]byte(text))
	return nil

	//drawingBytes := drawing.Data() // OfficeArtDgContainer within [MS-ODRAW] 2.2.13
	//from := 0
	//for {
	//	// instead of parsing binary PPT format, search text records directly
	//	pocketIdx := matchPocket(drawingBytes, from)
	//	if pocketIdx == -1 {
	//		break
	//	}
	//	// check if it is really a text record - recordType bytes must be preceded by 1-byte version and 3-byte instance
	//	// fields with zero values
	//	if pocketIdx >= 2 && bytes.Equal(drawingBytes[pocketIdx-headerRecordTypeOffset:pocketIdx], []byte{0x00, 0x00}) {
	//		var rec Record
	//		if drawingBytes[pocketIdx] == recordTypeTextBytesAtom.LowerPart() {
	//			rec, err = readRecord(drawing, int64(pocketIdx-headerRecordTypeOffset), recordTypeTextBytesAtom)
	//			if err != nil {
	//				return err
	//			}
	//			err = readTextFromTextBytesAtom(rec, out, utf16Decoder)
	//		} else {
	//			rec, err = readRecord(drawing, int64(pocketIdx-headerRecordTypeOffset), recordTypeTextCharsAtom)
	//			if err != nil {
	//				return err
	//			}
	//			err = readTextFromTextCharsAtom(rec, out, utf16Decoder)
	//		}
	//		if err != nil {
	//			return err
	//		}
	//	}
	//	from = pocketIdx + 2
	//}
	//return nil
}

func matchPocket(data []byte, from int) int {
	data = data[from:]
	n := len(data)
	for i := 0; i < n; i++ {
		switch data[i] {
		case recordTypeTextCharsAtom.LowerPart(), recordTypeTextBytesAtom.LowerPart():
			if i < n-1 && data[i+1] == 0x0F {
				return i + from
			}
		}
	}
	return -1
}

// readTextFromTextCharsAtom simply transforms UTF-16LE data into UTF-8 data
func readTextFromTextCharsAtom(atom Record, out *strings.Builder, dec *encoding.Decoder) error {
	dec.Reset()
	transformed, err := dec.Bytes(atom.Data())
	if err != nil {
		return err
	}
	out.Write(transformed)
	out.WriteByte(' ')
	return nil
}

func readTextFromTextBytesAtom(atom Record, out *strings.Builder, dec *encoding.Decoder) error {
	dec.Reset()
	transformed, err := decodeTextBytesAtom(atom.Data(), dec)
	if err != nil {
		return err
	}
	out.Write(transformed)
	out.WriteByte(' ')
	return nil
}

// decodeTextBytesAtom transforms text from TextBytesAtom, which is an array of bytes representing lower parts of UTF-16
// characters into UTF-8 data
func decodeTextBytesAtom(data []byte, dec *encoding.Decoder) ([]byte, error) {
	var (
		// buffer for UTF-16 char
		buf [2]byte
		err error
	)
	result := make([]byte, 0, len(data))
	for i := range data {
		// filling upper part of character with zero
		clear(buf[:])
		// fill lower part with byte
		buf[0] = data[i]

		// transform single UTF-16 char into UTF-8 rune and append it into result
		result, _, err = transform.Append(dec, result, buf[:])
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// skipRecords reads headers and skips data of records of provided types
func skipRecords(r io.ReaderAt, initialOffset int64, skippedRecordsTypes []RecordType) (int64, error) {
	offset := initialOffset

	for i := range skippedRecordsTypes {
		rec, err := readRecordHeaderOnly(r, offset, skippedRecordsTypes[i])
		if err != nil {
			if errors.Is(err, errMismatchRecordType) {
				continue
			}
			return 0, err
		}
		offset += int64(rec.Length() + headerSize)
	}

	return offset, nil
}

func readTargetRecordInstance(r Record, offset int64, instanceType InstanceType) (Record, error) {
	dataLen := int64(r.Length())
	if dataLen <= offset {
		return Record{}, errors.New("offset out of range")
	}
	for {
		rec, err := readRecordHeaderOnly(r, offset, recordTypeUnspecified)
		if err != nil {
			return Record{}, err
		}
		//fmt.Printf("target record type=%X instance=%X, current record type=%X instance=%X\n",
		//	instanceType.Type, instanceType.Instance, rec.Type(), rec.Instance())
		if rec.Type() == instanceType.Type && rec.Instance() == instanceType.Instance {
			rec.RecordData = make([]byte, rec.Length())
			_, err = r.ReadAt(rec.RecordData, offset+headerSize)
			if err != nil {
				fmt.Printf("read record body error: %v", err)
				return Record{}, err
			}
			return rec, nil
		}
		offset += int64(rec.Length() + headerSize)
		if offset >= int64(r.Length()) {
			return Record{}, errors.New("offset out of data range")
		}
	}
}
