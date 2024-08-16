package bitutil

func writeBitsPrimitive(objectToWrite *[16]byte, writeBitLength uint8, targetArray *[]byte, initalByteOffset uint64, initalBitOffset uint8) {
	target, source := *targetArray, *objectToWrite

	writedBits := uint8(0)
	writedBitOffset := uint8(0)
	writedByteOffset := uint64(0)

	for writedBits < writeBitLength {

		if initalBitOffset == 8 {
			initalByteOffset++
			initalBitOffset = 0
		}
		if writedBitOffset == 8 {
			writedByteOffset++
			writedBitOffset = 0
		}

		bit := (source[writedByteOffset] >> writedBitOffset) & 1

		target[initalByteOffset] &^= (1 << initalBitOffset)
		target[initalByteOffset] |= (bit << initalBitOffset)

		writedBitOffset++
		writedBits++
		initalBitOffset++
	}
}

func readBitsPrimitive(objectToRead *[16]byte, readBitLength uint8, arrayToRead *[]byte, initalByteOffset uint64, initalBitOffset uint8) {
	source, target := *arrayToRead, *objectToRead

	readBits := uint8(0)
	readBitOffset := uint8(0)
	readByteOffset := uint64(0)

	for readBits < readBitLength {

		if initalBitOffset == 8 {
			initalByteOffset++
			initalBitOffset = 0
		}
		if readBitOffset == 8 {
			readByteOffset++
			readBitOffset = 0
		}

		bit := (source[initalByteOffset] >> initalBitOffset) & 1

		target[readByteOffset] &^= (1 << readBitOffset)
		target[readByteOffset] |= (bit << readBitOffset)

		readBitOffset++
		readBits++
		initalBitOffset++
	}
}
