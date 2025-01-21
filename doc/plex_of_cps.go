package doc

import (
	"github.com/zhufuyin/mscfb/global"
	"io"
)

type PlexOfCps struct {
	iMac     int
	cbStruct int
	props    []*PropertyNode
}

func NewPlexOfCps(tblStream io.ReaderAt, start, cb, cbStruct int) (*PlexOfCps, error) {
	p := &PlexOfCps{}
	p.iMac = (cb - 4) / (4 + cbStruct)
	p.cbStruct = cbStruct
	for i := 0; i < p.iMac; i++ {
		propertyNode, err := p.BuildProperty(i, start, tblStream)
		if err != nil {
			return nil, err
		}
		p.props = append(p.props, propertyNode)
	}
	return p, nil
}

func (p *PlexOfCps) BuildProperty(index, offset int, tblStream io.ReaderAt) (*PropertyNode, error) {
	start, err := global.ReadUint32At(tblStream, int64(offset+4*index))
	if err != nil {
		return nil, err
	}
	end, err := global.ReadUint32At(tblStream, int64(offset+4*(index+1)))
	if err != nil {
		return nil, err
	}
	data := make([]byte, p.cbStruct)
	dataOffset := (4 * (p.iMac + 1)) + (p.cbStruct * index)
	_, err = tblStream.ReadAt(data, int64(offset+dataOffset))
	if err != nil {
		return nil, err
	}
	propertyNode := NewPropertyNode(int(start), int(end), data)
	return propertyNode, nil
}
