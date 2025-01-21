package ppt

import (
	"context"
	"errors"
	"github.com/zhufuyin/mscfb/cfb"
	"github.com/zhufuyin/mscfb/global"
	"io"
	"strings"
)

const (
	CurrentUserStreamName        = "Current User"
	PowerPointDocumentStreamName = "PowerPoint Document"
)

type PptFile struct {
	currentUserAtom    *CurrentUserAtom
	userEditAtoms      []*UserEditAtom
	latestUserEditAtom *UserEditAtom
	powerPointStream   *PowerPointStream
	//persistDirectoryOffsets []uint32
	//persistIdOffset map[uint32]int64
}

func NewPptFile(file io.Reader) (*PptFile, error) {
	ra := global.NewReaderAt(file)
	cfbFile, err := cfb.New(ra)
	if err != nil {
		return nil, err
	}
	ppt := &PptFile{}
	var currentUserStream, pptDocumentStream *cfb.File
	for _, stream := range cfbFile.File {
		if stream.Name == CurrentUserStreamName {
			currentUserStream = stream
		} else if stream.Name == PowerPointDocumentStreamName {
			pptDocumentStream = stream
		}
	}
	if currentUserStream == nil || pptDocumentStream == nil {
		return nil, errors.New("invalid ppt file")
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
	//offset := int64(currentUserAtom.offsetToCurrentEdit)
	//persistDirectoryOffsets := []uint32{}
	//userEditAtom, err := readUserEditAtom(pptDocumentStream, offset)
	persistDirectoryOffsets, err := ppt.readAllUserEditAtoms(pptDocumentStream)
	if err != nil {
		return nil, err
	}
	//ppt.latestUserEditAtom = userEditAtom
	ppt.powerPointStream = newPowerPointStream(pptDocumentStream, persistDirectoryOffsets)

	return ppt, nil
}

func (ppt *PptFile) readAllUserEditAtoms(pptDocumentStream io.ReaderAt) ([]uint32, error) {
	offset := int64(ppt.currentUserAtom.offsetToCurrentEdit)
	persistDirectoryOffsets := []uint32{}
	for {
		userEditAtom, err := readUserEditAtom(pptDocumentStream, offset)
		if err != nil {
			return nil, err
		}
		ppt.userEditAtoms = append(ppt.userEditAtoms, userEditAtom)
		persistDirectoryOffsets = append(persistDirectoryOffsets, userEditAtom.offsetPersistDirectory)
		offset = int64(userEditAtom.offsetLastEdit)
		if offset == 0 {
			break
		}
	}
	if len(ppt.userEditAtoms) > 0 {
		ppt.latestUserEditAtom = ppt.userEditAtoms[0]
	}

	return persistDirectoryOffsets, nil
}

func readUserEditAtom(pptDocumentStream io.ReaderAt, offset int64) (*UserEditAtom, error) {
	userEditAtomRecord, err := readRecord(pptDocumentStream, offset, recordTypeUserEditAtom)
	if err != nil {
		return nil, err
	}
	userEditAtom := &UserEditAtom{
		Record: userEditAtomRecord,
	}
	userEditAtom.parse()
	return userEditAtom, nil
}

func (ppt *PptFile) ExtractText(ctx context.Context) (string, error) {
	if ppt.latestUserEditAtom == nil {
		return "", nil
	}
	pptStream := ppt.powerPointStream
	err := pptStream.parsePersistDirectoryAtom()
	if err != nil {
		return "", err
	}
	//for _, userEditAtom := range ppt.userEditAtoms {
	//	err = pptStream.readDocumentContainer(ctx, userEditAtom)
	//	if err != nil {
	//		return "", err
	//	}
	//}
	err = pptStream.readDocumentContainer(ctx, ppt.latestUserEditAtom)
	if err != nil {
		return "", err
	}
	rawTexts, err := pptStream.extractText()
	if err != nil {
		return "", err
	}
	var targetTexts []string
	for _, rawText := range rawTexts {
		if len(rawText) > 0 && rawText != "*" {
			targetTexts = append(targetTexts, rawText)
		}
	}

	return strings.Join(targetTexts, "\n"), nil
}
