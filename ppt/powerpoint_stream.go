package ppt

import (
	"context"
	"fmt"
	"github.com/zhufuyin/mscfb/cfb"
)

type PowerPointStream struct {
	pptDocument             *cfb.File
	persistDirectoryAtoms   []*PersistDirectoryAtom // referred by offsetPersistDirectory in UserEditAtom
	persistDirectoryOffsets []uint32                // offsetPersistDirectory in UserEditAtom
	// persist object's offset in powerpoint document stream
	// value: persistId
	offset2PersistId    map[int64]uint32
	documentContainer   *DocumentContainer
	persistIdObjOffsets map[uint32]int64 // todo delete
}

func newPowerPointStream(pptDocument *cfb.File, persistDirectoryOffsets []uint32) *PowerPointStream {
	return &PowerPointStream{
		pptDocument:             pptDocument,
		persistDirectoryOffsets: persistDirectoryOffsets,
	}
}

// get all persistId and offset
func (s *PowerPointStream) parsePersistDirectoryAtom() error {
	s.offset2PersistId = make(map[int64]uint32)
	s.persistIdObjOffsets = make(map[uint32]int64)
	for _, offsetPersistDirectory := range s.persistDirectoryOffsets {
		offset := int64(offsetPersistDirectory)
		persistDirectoryAtomRecord, err := readRecord(s.pptDocument, offset, recordTypePersistDirectoryAtom)
		if err != nil {
			return err
		}
		persistDirectoryAtom := &PersistDirectoryAtom{
			Record: persistDirectoryAtomRecord,
		}
		err = persistDirectoryAtom.parse()
		if err != nil {
			return err
		}
		s.persistDirectoryAtoms = append(s.persistDirectoryAtoms, persistDirectoryAtom)
		for _, entry := range persistDirectoryAtom.rgPersistDirEntry {
			for persistId, blockOffset := range entry.persistIdOffset {
				s.offset2PersistId[blockOffset] = persistId
				if _, ok := s.persistIdObjOffsets[persistId]; !ok {
					s.persistIdObjOffsets[persistId] = blockOffset
				}
			}
		}
	}

	return nil
}

// read DocumentContainer whose offset is referred by docPersistIdRef in UserEditAtom
func (s *PowerPointStream) readDocumentContainer(ctx context.Context, userEditAtom *UserEditAtom) error {
	offset, ok := s.persistIdObjOffsets[userEditAtom.docPersistIdRef]
	if !ok {
		return fmt.Errorf("persistIdRef not found")
	}
	documentContainerRecord, err := readRecord(s.pptDocument, offset, recordTypeDocument)
	if err != nil {
		return err
	}
	documentContainer := &DocumentContainer{
		Record:       documentContainerRecord,
		pptDocStream: s,
	}
	err = documentContainer.parse(ctx)
	if err != nil {
		return err
	}
	s.documentContainer = documentContainer
	//fmt.Printf("documentContainer offset: %d\n", documentContainerRecord.offset)
	return nil
}

func (s *PowerPointStream) extractText() ([]string, error) {
	var texts []string

	// extract textCharsAtom and textBytesAtom in SlideListWithText
	directTexts, err := s.documentContainer.extractText()
	if err != nil {
		return nil, err
	}
	if len(directTexts) > 0 {
		texts = append(texts, directTexts...)
	}
	// slide list
	//slideTexts, err := s.extractTextFromSlides()
	//if err != nil {
	//	return nil, err
	//}
	//if len(slideTexts) > 0 {
	//	texts = append(texts, slideTexts...)
	//}

	// notes list
	//notesTexts, err := s.extractTextFromNotes()
	//if err != nil {
	//	return nil, err
	//}
	//if len(notesTexts) > 0 {
	//	texts = append(texts, notesTexts...)
	//}

	// master list
	//masterListTexts, err := s.extractTextFromMasterList()
	//if err != nil {
	//	return nil, err
	//}
	//if len(masterListTexts) > 0 {
	//	texts = append(texts, masterListTexts...)
	//}

	return texts, nil
}

func (s *PowerPointStream) extractTextFromSlides() ([]string, error) {
	var texts []string
	slideList := s.documentContainer.slideList
	if slideList == nil {
		return nil, nil
	}

	for _, slidePersistAtom := range slideList.slidePersistAtoms {
		slideTexts, err := s.extractTextFromSlide(slidePersistAtom)
		if err != nil {
			return nil, err
		}
		texts = append(texts, slideTexts...)
	}
	return texts, nil
}

func (s *PowerPointStream) extractTextFromSlide(slidePersistAtom *SlidePersistAtom) ([]string, error) {
	var texts []string
	persistIdRef := slidePersistAtom.persistIdRef
	offset, ok := s.persistIdObjOffsets[persistIdRef]
	if !ok {
		fmt.Printf("persistIdRef:%d not found\n", persistIdRef)
		return nil, fmt.Errorf("persistIdRef not found in persistIdObjOffsets")
	}
	slideRecord, err := readRecord(s.pptDocument, offset, recordTypeSlide)
	if err != nil {
		return nil, err
	}
	slide := &SlideContainer{
		Record: slideRecord,
	}
	err = slide.parse()
	if err != nil {
		return nil, err
	}
	slideTexts := slide.extractText()
	if len(slideTexts) > 0 {
		texts = append(texts, slideTexts...)
	}
	return texts, nil
}

func (s *PowerPointStream) extractTextFromNotes(documentContainer *DocumentContainer) ([]string, error) {
	var texts []string
	notesList := documentContainer.notesList
	if notesList == nil {
		return nil, nil
	}

	for _, notesPersistAtom := range notesList.rgNotesPersistAtom {
		persistIdRef := notesPersistAtom.persistIdRef
		offset, ok := s.persistIdObjOffsets[persistIdRef]
		if !ok {
			fmt.Printf("persistIdRef:%d not found\n", persistIdRef)
			return nil, fmt.Errorf("persistIdRef not found in persistIdObjOffsets")
		}
		notesRecord, err := readRecord(s.pptDocument, offset, recordTypeNotes)
		if err != nil {
			return nil, err
		}
		notes := &NotesContainer{
			Record: notesRecord,
		}
		err = notes.parse()
		if err != nil {
			return nil, err
		}
		notesTexts := notes.extractText()
		if len(notesTexts) > 0 {
			texts = append(texts, notesTexts...)
		}
	}
	return texts, nil
}

func (s *PowerPointStream) extractTextFromMasterList(documentContainer *DocumentContainer) ([]string, error) {
	var texts []string
	masterList := documentContainer.masterList
	if masterList == nil {
		return nil, nil
	}

	for _, masterPersistAtom := range masterList.rgMasterPersistAtom {
		persistIdRef := masterPersistAtom.persistIdRef
		offset, ok := s.persistIdObjOffsets[persistIdRef]
		if !ok {
			fmt.Printf("persistIdRef:%d not found\n", persistIdRef)
			return nil, fmt.Errorf("persistIdRef not found in persistIdObjOffsets")
		}
		masterRecord, err := readRecord(s.pptDocument, offset, recordTypeUnspecified)
		if err != nil {
			return nil, err
		}
		switch masterRecord.RecType {
		case recordTypeSlide:
			slideRecord := &SlideContainer{
				Record: masterRecord,
			}
			err = slideRecord.parse()
			if err != nil {
				return nil, err
			}
			slideTexts := slideRecord.extractText()
			if len(slideTexts) > 0 {
				texts = append(texts, slideTexts...)
			}
		case recordTypeMainMaster:
			mainMasterRecord := &MainMasterContainer{
				Record: masterRecord,
			}
			err = mainMasterRecord.parse()
			if err != nil {
				return nil, err
			}
			mainMasterTexts := mainMasterRecord.extractText()
			if len(mainMasterTexts) > 0 {
				texts = append(texts, mainMasterTexts...)
			}
		default:
			fmt.Printf("Unknown master record type %d\n", masterRecord.RecType)
		}
	}

	return texts, nil
}
