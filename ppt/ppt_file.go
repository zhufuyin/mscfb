package ppt

import "github.com/zhufuyin/mscfb/cfb"

type PptFile struct {
	currentUserAtom         CurrentUserAtom
	powerPointStream        *cfb.File
	persistDirectoryOffsets []uint32
}

func NewPptFile(file *cfb.Reader) (*PptFile, error) {
	ppt := &PptFile{}
	var currentUserStream *cfb.File
	for _, f := range file.File {
		switch f.Name {
		case "Current User":
			currentUserStream = f
		case "PowerPoint Document":
			ppt.powerPointStream = f
		}
	}
	err := isValidPPT(currentUserStream, ppt.powerPointStream)
	if err != nil {
		return nil, err
	}
	// parse current user stream
	currentUserStreamRecord, err := readRecord(currentUserStream, 0, recordTypeCurrentUserAtom)
	currentUserAtom := CurrentUserAtom{
		Record: currentUserStreamRecord,
	}
	currentUserAtom.parse()
	ppt.currentUserAtom = currentUserAtom

	/*
		read UserEditAtom and persist directory offset from powerpoint document stream,
		a chain linked by field [offsetLastEdit]
	*/
	offset := int64(currentUserAtom.offsetToCurrentEdit)
	for {
		userEditAtomRecord, err := readRecord(ppt.powerPointStream, offset, recordTypeUserEditAtom)
		if err != nil {
			return nil, err
		}
		userEditAtom := UserEditAtom{
			Record: userEditAtomRecord,
		}
		userEditAtom.parse()
		ppt.persistDirectoryOffsets = append(ppt.persistDirectoryOffsets, userEditAtom.offsetPersistDirectory)
		offset = int64(userEditAtom.offsetLastEdit)
		if offset == 0 {
			break
		}
	}
	// read persistDirectoryEntries
	for _, offset = range ppt.persistDirectoryOffsets {
		rgPersistDirEntry, err := readRecord(ppt.powerPointStream, offset, recordTypePersistDirectoryAtom)
		if err != nil {
			return nil, err
		}
		readPersistDirectoryEntry()

	}

	return ppt, nil
}
