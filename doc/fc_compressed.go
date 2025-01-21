package doc

// [MS-DOC] section 2.9.73
type FcCompressed struct {
	fc          uint32 // 30-bit
	fCompressed uint8  // 1-bit
	r1          uint8  // must be 0
}
