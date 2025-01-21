package doc

// describe CP
type PropertyNode struct {
	cpStart int
	cpEnd   int
	buf     []byte
}

func NewPropertyNode(cpStart, cpEnd int, buf []byte) *PropertyNode {
	pn := PropertyNode{}
	pn.cpStart = cpStart
	pn.cpEnd = cpEnd
	pn.buf = buf
	if cpStart < 0 {
		pn.cpStart = 0
	}
	if cpEnd < cpStart {
		cpEnd = cpStart
	}
	return &pn
}
