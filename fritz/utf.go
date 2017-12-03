package fritz

import (
	"unicode/utf16"
	"unicode/utf8"
)

func utf8To16LE(p []byte) []byte {
	bs := make([]byte, 0, 2*len(p))
	pos := 0
	for pos < len(p) {
		bytes, size := consumeNextRune(p[pos:])
		pos += size
		bs = append(bs, bytes...)
	}
	return bs
}

func consumeNextRune(p []byte) ([]byte, int) {
	r, size := utf8.DecodeRune(p)
	if r <= 0xffff {
		return []byte{uint8(r), uint8(r >> 8)}, size
	}
	r1, r2 := utf16.EncodeRune(r)
	return []byte{uint8(r1), uint8(r1 >> 8), uint8(r2), uint8(r2 >> 8)}, size
}
