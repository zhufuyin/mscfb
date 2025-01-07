package ppt

type SlidePersistAtom struct {
	Record
	persistIdRef uint32 // offset of the SlideContainer record
	slideId      uint32 // identifier for a presentation slide
}

func (s *SlidePersistAtom) parse() error {
	s.persistIdRef = s.Uint32At(0)
	s.slideId = s.Uint32At(12)
	return nil
}
