package fritz

import "strconv"

type bitMasked struct {
	Functionbitmask string
}

func (b bitMasked) hasMask(mask int64) bool {
	bitMask, err := strconv.ParseInt(b.Functionbitmask, 10, 64)
	if err != nil {
		return false
	}
	return (bitMask & mask) != 0
}
