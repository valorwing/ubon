package uIntEncoderHelper

import (
	"errors"
	bitcode "ubon/internal/bitCode"
	"ubon/internal/ubonOpCodes"
	"ubon/internal/writeOnlyBitStream"
	"unsafe"
)

var rawUIntSize = unsafe.Sizeof(uint(0))

const (
	uInt32ByteSize = 4
	uInt64ByteSize = 8
)

func UntypedUIntBitLenUnsafe() uint8 {
	if rawUIntSize == uInt32ByteSize {
		return 32
	} else {
		return 64
	}
}

func UntypedUIntBitLen() (uint8, error) {
	if rawUIntSize == uInt32ByteSize {
		return 32, nil
	} else if rawUIntSize == uInt64ByteSize {
		return 64, nil
	}
	return 0, errors.New("untypedUIntBitLen unexpected size of int (32/64 allowed)")
}

// Auto Write Allow UInt Bit Trim
func WriteAutoUnsignedIntToBitStream(value any, ws *writeOnlyBitStream.WriteOnlyBitStream) error {
	valueToCheckType := value
	switch valueToCheckType.(type) {
	case uint8:
		return WriteUInt8ToBitStream(valueToCheckType.(uint8), ws)
	case uint16:
		value16 := valueToCheckType.(uint16)
		if value16 == uint16(uint8(value16)) {
			return WriteUInt8ToBitStream(uint8(value16), ws)
		} else {
			return WriteUInt16ToBitStream(value16, ws)
		}
	case uint32:
		value32 := uint32(valueToCheckType.(uint32))
		if value32 == uint32(uint8(value32)) {
			return WriteUInt8ToBitStream(uint8(value32), ws)
		} else if value32 == uint32(uint16(value32)) {
			return WriteUInt16ToBitStream(uint16(value32), ws)
		} else {
			return WriteUInt32ToBitStream(value32, ws)
		}
	case uint64:
		value64 := uint64(valueToCheckType.(uint64))
		if value64 == uint64(uint8(value64)) {
			return WriteUInt8ToBitStream(uint8(value64), ws)
		} else if value64 == uint64(uint16(value64)) {
			return WriteUInt16ToBitStream(uint16(value64), ws)
		} else if value64 == uint64(uint32(value64)) {
			return WriteUInt32ToBitStream(uint32(value64), ws)
		} else {
			return WriteUInt64ToBitStream(value64, ws)
		}
	case uint:
		if rawUIntSize == uInt32ByteSize {
			value32 := uint32(valueToCheckType.(uint))
			if value32 == uint32(uint8(value32)) {
				return WriteUInt8ToBitStream(uint8(value32), ws)
			} else if value32 == uint32(int16(value32)) {
				return WriteUInt16ToBitStream(uint16(value32), ws)
			} else {
				return WriteUInt32ToBitStream(value32, ws)
			}
		} else if rawUIntSize == uInt64ByteSize {
			value64 := uint64(valueToCheckType.(uint))
			if value64 == uint64(uint8(value64)) {
				return WriteUInt8ToBitStream(uint8(value64), ws)
			} else if value64 == uint64(uint16(value64)) {
				return WriteUInt16ToBitStream(uint16(value64), ws)
			} else if value64 == uint64(int32(value64)) {
				return WriteUInt32ToBitStream(uint32(value64), ws)
			} else {
				return WriteUInt64ToBitStream(value64, ws)
			}
		} else {
			return errors.New("writeAutoUnsignedIntToBitStream unprefixed int error size detection")
		}
	default:
		return errors.New("writeAutoUnsignedIntToBitStream only accepts int types (int8/16/32/64 and int)")
	}
}

func WriteUInt8ToBitStream(value uint8, ws *writeOnlyBitStream.WriteOnlyBitStream) error {
	if ws == nil {
		return errors.New("ws aka writeOnlyBitStream can't be nil")
	}
	//opCode for safe clone
	bitCodeToWrite := ubonOpCodes.UBON_OP_NEXT_UNSIGNED_INT.Clone()
	//mask 0 - 8 bits
	bitCodeToWrite.Append(false)
	//write single int8 byte
	bitCodeToWrite.AppendBitCode(bitcode.NewBitCodeFromBytes(byte(value)))
	//write to stream
	ws.AppendBitCode(bitCodeToWrite)
	return nil
}

func WriteUInt16ToBitStream(value uint16, ws *writeOnlyBitStream.WriteOnlyBitStream) error {
	if ws == nil {
		return errors.New("ws aka writeOnlyBitStream can't be nil")
	}
	//opCode for safe clone
	bitCodeToWrite := ubonOpCodes.UBON_OP_NEXT_UNSIGNED_INT.Clone()
	//mask 10 - 16 bits
	bitCodeToWrite.AppendBitCode(bitcode.NewBitCodeWithBoolCode(true, false))
	//write 2 bytes little endian
	bitCodeToWrite.AppendBitCode(bitcode.NewBitCodeFromBytes(byte(value), byte(value>>8)))
	//write to stream
	ws.AppendBitCode(bitCodeToWrite)
	return nil
}

func WriteUInt32ToBitStream(value uint32, ws *writeOnlyBitStream.WriteOnlyBitStream) error {
	if ws == nil {
		return errors.New("ws aka writeOnlyBitStream can't be nil")
	}
	//opCode for safe clone
	bitCodeToWrite := ubonOpCodes.UBON_OP_NEXT_UNSIGNED_INT.Clone()
	//mask 110 - 32 bits
	bitCodeToWrite.AppendBitCode(bitcode.NewBitCodeWithBoolCode(true, true, false))
	//write 4 bytes little endian
	bitCodeToWrite.AppendBitCode(bitcode.NewBitCodeFromBytes(
		byte(value), byte(value>>8),
		byte(value>>16), byte(value>>24)))
	//write to stream
	ws.AppendBitCode(bitCodeToWrite)
	return nil
}

func WriteUInt64ToBitStream(value uint64, ws *writeOnlyBitStream.WriteOnlyBitStream) error {
	if ws == nil {
		return errors.New("ws aka writeOnlyBitStream can't be nil")
	}
	//opCode for safe clone
	bitCodeToWrite := ubonOpCodes.UBON_OP_NEXT_UNSIGNED_INT.Clone()
	//mask 111 - 64 bits (mask len 3 bits last value)
	bitCodeToWrite.AppendBitCode(bitcode.NewBitCodeWithBoolCode(true, true, true))
	//write 8 bytes little endian
	bitCodeToWrite.AppendBitCode(bitcode.NewBitCodeFromBytes(
		byte(value), byte(value>>8),
		byte(value>>16), byte(value>>24),
		byte(value>>32), byte(value>>40),
		byte(value>>48), byte(value>>56)))
	//write to stream
	ws.AppendBitCode(bitCodeToWrite)
	return nil
}
