package writeOnlyBitStream

import bitcode "ubon/internal/bitCode"

type WriteOnlyBitStream struct {
	data          []byte
	lastBitOffset int
}

func NewWriteOnlyBitStream() *WriteOnlyBitStream {
	return &WriteOnlyBitStream{
		data:          []byte{0},
		lastBitOffset: 0,
	}
}

func (bs *WriteOnlyBitStream) AppendBitCode(b bitcode.BitCode) {
	readBitCodeOffset := 0

	bitCodeValue := byte(0)
	readedBits := 0
	hasNext := true

	for hasNext {
		bitCodeValue, readedBits, hasNext = b.ReadSerialized(readBitCodeOffset, 8-bs.lastBitOffset)
		//write
		bitCodeValue >>= bs.lastBitOffset
		bs.data[len(bs.data)-1] |= bitCodeValue

		bs.lastBitOffset += readedBits
		if bs.lastBitOffset == 8 {
			bs.data = append(bs.data, 0)
			bs.lastBitOffset = 0
		}
		readBitCodeOffset += readedBits

	}
}

func (bs *WriteOnlyBitStream) Bytes() []byte {
	if bs.lastBitOffset == 0 {
		return bs.data[:len(bs.data)-1]
	}
	return bs.data
}
