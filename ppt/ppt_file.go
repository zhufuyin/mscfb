package ppt

import (
	"context"
	"fmt"
	"github.com/zhufuyin/mscfb/cfb"
	"io"
	"strings"
)

type PptFile struct {
	currentUserAtom  *CurrentUserAtom
	userEditAtoms    []UserEditAtom
	powerPointStream *PowerPointStream
	//persistDirectoryOffsets []uint32
	//persistIdOffset map[uint32]int64
}

func NewPptFile(file io.Reader) (*PptFile, error) {
	ra := ToReaderAt(file)
	cfbFile, err := cfb.New(ra)
	if err != nil {
		return nil, err
	}
	ppt := &PptFile{}
	var currentUserStream, pptDocumentStream *cfb.File
	for _, f := range cfbFile.File {
		switch f.Name {
		case "Current User":
			currentUserStream = f
		case "PowerPoint Document":
			pptDocumentStream = f
		default:
			fmt.Println("Unknown File: ", f.Name)
		}
	}
	err = isValidPPT(currentUserStream, pptDocumentStream)
	if err != nil {
		return nil, err
	}
	// parse current user stream
	currentUserStreamRecord, err := readRecord(currentUserStream, 0, recordTypeCurrentUserAtom)
	currentUserAtom := &CurrentUserAtom{
		Record: currentUserStreamRecord,
	}
	currentUserAtom.parse()
	ppt.currentUserAtom = currentUserAtom

	/*
		read UserEditAtom and persist directory offset from powerpoint document stream,
		UserEditAtom chain is linked by field [offsetLastEdit]
	*/
	offset := int64(currentUserAtom.offsetToCurrentEdit)
	persistDirectoryOffsets := []uint32{}
	for {
		userEditAtomRecord, err := readRecord(pptDocumentStream, offset, recordTypeUserEditAtom)
		if err != nil {
			return nil, err
		}
		userEditAtom := UserEditAtom{
			Record: userEditAtomRecord,
		}
		userEditAtom.parse()
		ppt.userEditAtoms = append(ppt.userEditAtoms, userEditAtom)
		persistDirectoryOffsets = append(persistDirectoryOffsets, userEditAtom.offsetPersistDirectory)
		offset = int64(userEditAtom.offsetLastEdit)
		if offset == 0 {
			break
		}
	}
	// read persistDirectoryEntries
	//for _, userEditAtom := range ppt.userEditAtoms {
	//	persistDirOffset := userEditAtom.offsetPersistDirectory
	//	rgPersistDirEntry, err := readRecord(pptDocumentStream, int64(persistDirOffset), recordTypePersistDirectoryAtom)
	//	if err != nil {
	//		return nil, err
	//	}
	//	persistDirectoryAtom := PersistDirectoryAtom{
	//		Record: rgPersistDirEntry,
	//	}
	//	err = persistDirectoryAtom.parse()
	//	if err != nil {
	//		return nil, err
	//	}
	//	// put all persistId and object offset into ppt
	//	for _, entry := range persistDirectoryAtom.rgPersistDirEntry {
	//		for persistId, objOffset := range entry.persistIdOffset {
	//			ppt.persistIdOffset[persistId] = objOffset
	//		}
	//	}
	//}
	ppt.powerPointStream = newPowerPointStream(pptDocumentStream, persistDirectoryOffsets)

	return ppt, nil
}

func (ppt *PptFile) extractText(ctx context.Context) (string, error) {
	pptStream := ppt.powerPointStream
	err := pptStream.parsePersistDirectoryAtom()
	if err != nil {
		return "", err
	}
	var texts []string
	for _, userEditAtom := range ppt.userEditAtoms {
		err = pptStream.readDocumentContainer(ctx, userEditAtom)
		if err != nil {
			return "", err
		}
		userEditTexts, err := pptStream.extractText()
		if err != nil {
			return "", err
		}
		if len(userEditTexts) > 0 {
			texts = append(texts, userEditTexts...)
		}
	}

	return strings.Join(texts, "\n"), nil
}
