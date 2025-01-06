package ppt

import "github.com/zhufuyin/mscfb/cfb"

type PowerPointStream struct {
	Record
	PptDocument *cfb.File
}
