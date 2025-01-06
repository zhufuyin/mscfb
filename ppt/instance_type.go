package ppt

type InstanceType struct {
	Version  uint16
	Instance uint16
	Type     RecordType
}

var (
	MasterListWithTextContainerType = InstanceType{
		Version:  0xF,
		Instance: 0x001,
		Type:     recordTypeSlideListWithText,
	}
	SlideListWithTextContainerType = InstanceType{
		Version:  0xF,
		Instance: 0x000,
		Type:     recordTypeSlideListWithText,
	}
	NotesListWithTextContainerType = InstanceType{
		Version:  0xF,
		Instance: 0x002,
		Type:     recordTypeSlideListWithText,
	}
)
