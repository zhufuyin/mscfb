package ppt

import "github.com/zhufuyin/mscfb/cfb"

type PersistDirectory struct {
	Document        *cfb.File
	PersistIdOffset map[uint32]int64
}

func NewPersistDirectory(document *cfb.File) *PersistDirectory {
	return &PersistDirectory{
		Document: document,
	}
}

func (p *PersistDirectory) Init() {

}
