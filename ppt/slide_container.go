package ppt

// SlideContainer [MS-PPT] 2.5.1
type SlideContainer struct {
	Record
	drawing *DrawingContainer
}

func (s *SlideContainer) parse() error {
	offset := int64(0)
	drawingType := InstanceType{
		Type:     recordTypeDrawing,
		Instance: 0x000,
	}
	drawingRecord, err := readTargetRecordInstance(s.Record, offset, drawingType)
	if err != nil {
		return err
	}
	s.drawing = &DrawingContainer{
		Record: drawingRecord,
	}
	err = s.drawing.parse()
	if err != nil {
		return err
	}
	return nil
}

func (s *SlideContainer) extractText() []string {
	return s.drawing.extractText()
}
