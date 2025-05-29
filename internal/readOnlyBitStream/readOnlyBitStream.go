package readOnlyBitStream

import (
	"errors"
	bitcode "ubon/internal/bitCode"
)

type ReadOnlyBitStream struct {
	data           []byte
	frontBitOffset uint8
}

func NewReadOnlyBitStream(input []byte) ReadOnlyBitStream {
	return ReadOnlyBitStream{data: input, frontBitOffset: 0}
}

func (bs *ReadOnlyBitStream) ReadBitCode(length int) (*bitcode.BitCode, error) {
	retVal := bitcode.NewZeroBitCodeWithLength(length)

	for i := 0; i < length; i++ {
		if len(bs.data) == 0 {
			return nil, errors.New("unexpected end stream")
		}
		if (bs.data[0] & (1 << (7 - bs.frontBitOffset))) != 0 {
			retVal.SetBit(i)
		}
		bs.frontBitOffset++
		if bs.frontBitOffset >= 8 {
			bs.frontBitOffset = 0
			if len(bs.data) <= 1 {
				bs.data = bs.data[0:]
			} else {
				bs.data = bs.data[1:]
			}
		}
	}
	return &retVal, nil
}
