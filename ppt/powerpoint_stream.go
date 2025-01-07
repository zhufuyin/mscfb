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
	persistIdObjOffsets     map[uint32]int64        // persistId -> persist object's offset in powerpoint document stream
	documentContainer       *DocumentContainer
}

func newPowerPointStream(pptDocument *cfb.File, persistDirectoryOffsets []uint32) *PowerPointStream {
	return &PowerPointStream{
		pptDocument:             pptDocument,
		persistDirectoryOffsets: persistDirectoryOffsets,
	}
}

func (s *PowerPointStream) parsePersistDirectoryAtom() error {
	s.persistIdObjOffsets = make(map[uint32]int64)
	for _, offsetPersistDirectory := range s.persistDirectoryOffsets {
		persistDirectoryAtomRecord, err := readRecord(s.pptDocument, int64(offsetPersistDirectory), recordTypePersistDirectoryAtom)
		if err != nil {
			return err
		}
		persistDirectoryAtom := PersistDirectoryAtom{
			Record: persistDirectoryAtomRecord,
		}
		err = persistDirectoryAtom.parse()
		if err != nil {
			return err
		}
		s.persistDirectoryAtoms = append(s.persistDirectoryAtoms, &persistDirectoryAtom)
		for _, entry := range persistDirectoryAtom.rgPersistDirEntry {
			for persistId, blockOffset := range entry.persistIdOffset {
				s.persistIdObjOffsets[persistId] = blockOffset
			}
		}
	}
	return nil
}

// read DocumentContainer whose offset is referred by docPersistIdRef in UserEditAtom
func (s *PowerPointStream) readDocumentContainer(ctx context.Context, userEditAtom UserEditAtom) error {
	documentContainerRecord, err := readRecord(s.pptDocument, int64(userEditAtom.docPersistIdRef), recordTypeDocument)
	if err != nil {
		return err
	}
	s.documentContainer = &DocumentContainer{
		Record: documentContainerRecord,
	}
	err = s.documentContainer.parse(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *PowerPointStream) extractText() ([]string, error) {
	var texts []string
	documentContainer := s.documentContainer
	directTexts, err := documentContainer.extractText()
	if err != nil {
		return nil, err
	}
	if len(directTexts) > 0 {
		texts = append(texts, directTexts...)
	}
	// slide list
	slideList := documentContainer.slideList
	for _, slidePersistAtom := range slideList.slidePersistAtoms {
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
	}
	// notes list
	notesList := s.documentContainer.notesList
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
	// master list
	masterList := s.documentContainer.masterList
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
