package bitcode

import (
	"bytes"
)

type BitCode struct {
	bitPointer uint8
	storage    []byte
}

func NewZeroBitCodeWithLength(length int) BitCode {
	byteLen := length / 8
	bitLen := length % 8
	if bitLen == 0 && byteLen == 0 {
		byteLen = 1
	} else {
		byteLen = byteLen + 1
	}
	return BitCode{bitPointer: uint8(bitLen), storage: make([]byte, byteLen)}
}

func NewBitCodeWithBoolCode(code ...bool) BitCode {
	retVal := NewZeroBitCodeWithLength(len(code))
	for i, v := range code {
		if v {
			retVal.storage[i/8] |= 1 << (7 - uint(i%8))
		}
	}
	return retVal
}

func NewBitCodeFromBytes(bytesData ...byte) BitCode {
	retVal := NewZeroBitCodeWithLength(len(bytesData) * 8)
	bitPoiner := 0
	for _, v := range bytesData {
		for i := 0; i < 8; i++ {

			if (v & (1 << (7 - uint(i%8)))) != 0 {
				retVal.storage[bitPoiner/8] |= 1 << (7 - uint(bitPoiner%8))
			}

			bitPoiner++
		}
	}
	return retVal
}

func (b BitCode) BitLength() int {
	if len(b.storage) == 0 {
		return 0
	} else if len((b.storage)) == 1 {
		return int(b.bitPointer)
	}

	bitLen := (len(b.storage) - 1) * 8
	if b.bitPointer != 0 {
		bitLen += int(b.bitPointer)
	}

	return bitLen
}

func (b *BitCode) SetBit(bitIndex int) {
	b.storage[bitIndex/8] |= 1 << (7 - uint(bitIndex%8))
}

func (b BitCode) Clone() BitCode {
	retVal := BitCode{bitPointer: b.bitPointer, storage: make([]byte, len(b.storage))}
	copy(retVal.storage, b.storage)
	return retVal
}

func (b BitCode) CloneAppend(bit bool) BitCode {
	retVal := b.Clone()
	retVal.Append(bit)
	return retVal
}

func (b *BitCode) AppendBitCode(other BitCode) {
	//fast append
	oldLen := b.BitLength()
	bytesAppendCount := 0
	if len(other.storage) > 1 {
		bytesAppendCount = len(other.storage) - 1
	}
	appendSlice := make([]byte, 0)
	for i := 0; i < bytesAppendCount; i++ {
		appendSlice = append(appendSlice, 0)
	}
	b.storage = append(b.storage, appendSlice...)
	b.bitPointer += other.bitPointer
	if b.bitPointer >= 8 {
		b.bitPointer -= 8
		b.storage = append(b.storage, 0)
	}
	for i := 0; i < other.BitLength(); i++ {
		if other.GetBit(i) {
			b.SetBit(oldLen + i)
		}
	}
}

func (b *BitCode) Append(bit bool) {
	if bit {
		b.storage[len(b.storage)-1] |= 1 << (7 - uint(b.bitPointer%8))
	}
	b.bitPointer++
	if b.bitPointer >= 8 {
		b.bitPointer = 0
		b.storage = append(b.storage, 0)
	}
}

func (b BitCode) GetBit(i int) bool {
	bit := (b.storage[i/8] & (1 << (7 - uint(i%8)))) != 0
	return bit
}

// debug
func (b BitCode) String() string {
	//no lib
	retVal := ""
	for i := 0; i < b.BitLength(); i++ {
		s := "0"
		if (b.storage[i/8] & (1 << (7 - uint(i%8)))) != 0 {
			s = "1"
		}
		retVal = retVal + s
	}
	return retVal
}

func (b *BitCode) Clear() {
	b.storage = b.storage[:1]
	b.bitPointer = 0
	b.storage[0] = 0
}

func (b BitCode) Equal(other BitCode) bool {
	//i'm lazy yeah yeah yeah
	return b.bitPointer == other.bitPointer && bytes.Equal(b.storage, other.storage)
}

func (b BitCode) Bytes() []byte {
	if len(b.storage) == 0 {
		return nil
	}
	end := len(b.storage) - 1
	if b.bitPointer > 0 {
		end++
	}
	return b.storage[:end]
}

func (b BitCode) ReadSerialized(offset, maxReadLen int) (value byte, readedBits int, hasNext bool) {
	if maxReadLen-offset > 8 {
		//I know it's not good to panic but
		//this shouldn't happen in the program
		//it's an algorithmic error the best way
		// out of which is an emergency exit from
		//the program for feedback
		panic("a byte has a capacity of 8. number of bytes read is more than 8")
	}

	value = byte(0)
	//the return value always starts
	// writing from the most significant byte for convenience,
	// the byte is always represented as an entity that is filled from left to right
	valueIndex := 7
	i := offset
	readedBits = 0
	for ; i < offset+maxReadLen && i < b.BitLength(); i++ {

		if (b.storage[i/8] & (1 << (7 - uint(i%8)))) != 0 {
			value |= 1 << valueIndex
		}
		valueIndex--
		readedBits++
	}
	hasNext = i < b.BitLength()
	return
}

func (b BitCode) Hash() uint64 {
	var h uint64 = 146527
	const prime uint64 = 1099511628211

	fullBytes := len(b.storage) - 1
	if fullBytes < 0 {
		fullBytes = 0
	}

	for i := 0; i < fullBytes; i++ {
		h ^= uint64(b.storage[i])
		h *= prime
	}

	if len(b.storage) > 0 {
		lastIdx := len(b.storage) - 1
		if b.bitPointer == 0 {

			h ^= uint64(b.storage[lastIdx])
			h *= prime
		} else {
			mask := byte(0xFF) << (8 - b.bitPointer)
			masked := b.storage[lastIdx] & mask
			h ^= uint64(masked)
			h *= prime
			h ^= uint64(b.bitPointer)
			h *= prime
		}
	}

	return h
}
