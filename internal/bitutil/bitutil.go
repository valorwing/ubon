package bitutil

func SetBit(b byte, bitPosition uint8) byte {
	retVal := b
	switch bitPosition {
	case 0:
		retVal = b | 0b10000000
	case 1:
		retVal = b | 0b01000000
	case 2:
		retVal = b | 0b00100000
	case 3:
		retVal = b | 0b00010000
	case 4:
		retVal = b | 0b00001000
	case 5:
		retVal = b | 0b00000100
	case 6:
		retVal = b | 0b00000010
	case 7:
		retVal = b | 0b00000001
	default:
		panic("set invalid bit position")
	}
	return retVal
}
func ResetBit(b byte, bitPosition uint8) byte {
	retVal := b
	switch bitPosition {
	case 0:
		retVal = b & 0b01111111
	case 1:
		retVal = b & 0b10111111
	case 2:
		retVal = b & 0b11011111
	case 3:
		retVal = b & 0b11101111
	case 4:
		retVal = b & 0b11110111
	case 5:
		retVal = b & 0b11111011
	case 6:
		retVal = b & 0b11111101
	case 7:
		retVal = b & 0b11111110
	default:
		panic("set invalid bit position")
	}
	return retVal
}
func ReadBit(b byte, bitPosition uint8) bool {
	retVal := b
	switch bitPosition {
	case 0:
		retVal = b & 0b10000000
	case 1:
		retVal = b & 0b01000000
	case 2:
		retVal = b & 0b00100000
	case 3:
		retVal = b & 0b00010000
	case 4:
		retVal = b & 0b00001000
	case 5:
		retVal = b & 0b00000100
	case 6:
		retVal = b & 0b00000010
	case 7:
		retVal = b & 0b00000001
	default:
		panic("set invalid bit position")
	}
	return retVal != 0
}

func WriteBitsPrimitive(objectToWrite *[16]byte, writeBitLength uint8, targetArray *[]byte, initalByteOffset uint64, initalBitOffset uint8) {
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
		//Optimized bit read
		bit := (source[writedByteOffset] >> (7 - writedBitOffset)) & 1
		//Optimized bit reset
		target[initalByteOffset] &^= (1 << (7 - initalBitOffset))
		//Optimized bit write
		target[initalByteOffset] |= (bit << (7 - initalBitOffset))

		writedBitOffset++
		writedBits++
		initalBitOffset++
	}
}

func ReadBitsPrimitive(objectToRead *[16]byte, readBitLength uint8, arrayToRead *[]byte, initalByteOffset uint64, initalBitOffset uint8) {
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
		//Optimized bit read
		bit := (source[initalByteOffset] >> (7 - initalBitOffset)) & 1
		//Optimized bit reset
		target[readByteOffset] &^= 1 << (7 - readBitOffset)
		//Optimized bit write
		target[readByteOffset] |= (bit << (7 - readBitOffset))

		readBitOffset++
		readBits++
		initalBitOffset++
	}
}
