package doc

type Pcd struct {
	fNoParaLast uint8
	fR1         uint8
	fDirty      uint8
	fR2         uint16
	fc          *FcCompressed
	prm         *Prm
}
