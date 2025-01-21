package global

import (
	"encoding/binary"
	"errors"
	"golang.org/x/text/encoding/unicode"
	"io"
	"slices"
)

var (
	Utf16Decoder = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
)

type ReaderAt struct {
	r   io.Reader
	buf []byte
}

func NewReaderAt(r io.Reader) io.ReaderAt {
	ra, ok := r.(io.ReaderAt)
	if ok {
		return ra
	}
	return &ReaderAt{
		r: r,
	}
}

func (r *ReaderAt) ReadAt(p []byte, off int64) (n int, err error) {
	if int(off)+len(p) > len(r.buf) {
		err := r.grow(int(off) + len(p))
		if err != nil {
			return 0, err
		}
	}
	return ReadBytes(r.buf, p, off)
}

func (r *ReaderAt) grow(newSize int) error {
	if cap(r.buf) < newSize {
		r.buf = slices.Grow(r.buf, newSize-cap(r.buf))
	}

	newPart := r.buf[len(r.buf):newSize]
	n, err := r.r.Read(newPart)
	switch {
	case err == nil:
		r.buf = r.buf[:newSize]
	case errors.Is(err, io.EOF):
		r.buf = r.buf[:len(r.buf)+n]
	default:
		return err
	}
	return nil
}

func ReadBytes(src []byte, dst []byte, off int64) (n int, err error) {
	idx := 0
	for i := int(off); i < len(src) && idx < len(dst); i, idx = i+1, idx+1 {
		dst[idx] = src[i]
	}
	if idx != len(dst) {
		return idx, io.EOF
	}
	return idx, nil
}

func ReadUint32At(r io.ReaderAt, offset int64) (uint32, error) {
	buf := make([]byte, 4)
	_, err := r.ReadAt(buf, offset)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(buf), nil
}

func ReadUint64At(r io.ReaderAt, offset int64) (uint64, error) {
	buf := make([]byte, 8)
	_, err := r.ReadAt(buf, offset)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(buf), nil
}

func ReadUint16At(r io.ReaderAt, offset int64) (uint16, error) {
	buf := make([]byte, 2)
	_, err := r.ReadAt(buf, offset)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(buf), nil
}

func ReadByteAt(r io.ReaderAt, offset int64) (byte, error) {
	buf := make([]byte, 1)
	_, err := r.ReadAt(buf, offset)
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}
