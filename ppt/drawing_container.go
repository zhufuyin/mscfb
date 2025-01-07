package ppt

// DrawingContainer [MS-PPT] 2.5.13
type DrawingContainer struct {
	Record
	officeArgDg *OfficeArtDGContainer
}

func (d *DrawingContainer) parse() error {
	officeArgDgContainerRecord, err := readRecord(d, 0, recordTypeOfficeArtDgContainer)
	if err != nil {
		return err
	}
	officeArgDgContainer := &OfficeArtDGContainer{
		Record: officeArgDgContainerRecord,
	}
	err = officeArgDgContainer.parse()
	if err != nil {
		return err
	}
	d.officeArgDg = officeArgDgContainer
	return nil
}

func (d *DrawingContainer) extractText() []string {
	return d.officeArgDg.extractText()
}
