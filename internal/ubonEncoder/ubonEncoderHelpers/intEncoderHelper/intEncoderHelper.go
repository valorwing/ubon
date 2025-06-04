package intEncoderHelper

import (
	"errors"
	bitcode "ubon/internal/bitCode"
	uIntEncoderHelper "ubon/internal/ubonEncoder/ubonEncoderHelpers/uintEncoderHelper"

	"ubon/internal/ubonOpCodes"
	"ubon/internal/writeOnlyBitStream"
	"unsafe"
)

var rawIntSize = unsafe.Sizeof(int(0))

const (
	int32ByteSize = 4
	int64ByteSize = 8
)

func UntypedIntBitLenUnsafe() uint8 {
	if rawIntSize == int32ByteSize {
		return 32
	} else {
		return 64
	}
}

func UntypedIntBitLen() (uint8, error) {
	if rawIntSize == int32ByteSize {
		return 32, nil
	} else if rawIntSize == int64ByteSize {
		return 64, nil
	}
	return 0, errors.New("untypedIntBitLen unexpected size of int (32/64 allowed)")
}

// Auto Write Allow Int Bit Trim and Auto Unsign For compact
func WriteAutoIntToBitStream(value any, ws *writeOnlyBitStream.WriteOnlyBitStream) error {
	valueToCheckType := value

	switch valueToCheckType.(type) {
	case int8:
		return WriteInt8ToBitStream(valueToCheckType.(int8), ws)
	case int16:
		value16 := valueToCheckType.(int16)
		if (uint16(value16) & (1 << 15)) == 0 {
			//Unsigned
			return uIntEncoderHelper.WriteAutoUnsignedIntToBitStream(uint16(value16), ws)
		}
		if value16 == int16(int8(value16)) {
			return WriteInt8ToBitStream(int8(value16), ws)
		} else {
			return WriteInt16ToBitStream(value16, ws)
		}
	case int32:
		value32 := int32(valueToCheckType.(int32))
		if (uint32(value32) & (1 << 31)) == 0 {
			//Unsigned
			return uIntEncoderHelper.WriteAutoUnsignedIntToBitStream(uint32(value32), ws)
		}
		if value32 == int32(int8(value32)) {
			return WriteInt8ToBitStream(int8(value32), ws)
		} else if value32 == int32(int16(value32)) {
			return WriteInt16ToBitStream(int16(value32), ws)
		} else {
			return WriteInt32ToBitStream(value32, ws)
		}
	case int64:
		value64 := int64(valueToCheckType.(int64))
		if (uint64(value64) & (1 << 63)) == 0 {
			//Unsigned
			return uIntEncoderHelper.WriteAutoUnsignedIntToBitStream(uint64(value64), ws)
		}
		if value64 == int64(int8(value64)) {
			return WriteInt8ToBitStream(int8(value64), ws)
		} else if value64 == int64(int16(value64)) {
			return WriteInt16ToBitStream(int16(value64), ws)
		} else if value64 == int64(int32(value64)) {
			return WriteInt32ToBitStream(int32(value64), ws)
		} else {
			return WriteInt64ToBitStream(value64, ws)
		}
	case int:
		if rawIntSize == int32ByteSize {
			value32 := int32(valueToCheckType.(int))
			if (uint32(value32) & (1 << 31)) == 0 {
				//Unsigned
				return uIntEncoderHelper.WriteAutoUnsignedIntToBitStream(uint32(value32), ws)
			}
			if value32 == int32(int8(value32)) {
				return WriteInt8ToBitStream(int8(value32), ws)
			} else if value32 == int32(int16(value32)) {
				return WriteInt16ToBitStream(int16(value32), ws)
			} else {
				return WriteInt32ToBitStream(value32, ws)
			}
		} else if rawIntSize == int64ByteSize {
			value64 := int64(valueToCheckType.(int))
			if (uint64(value64) & (1 << 63)) == 0 {
				//Unsigned
				return uIntEncoderHelper.WriteAutoUnsignedIntToBitStream(uint64(value64), ws)
			}
			if value64 == int64(int8(value64)) {
				return WriteInt8ToBitStream(int8(value64), ws)
			} else if value64 == int64(int16(value64)) {
				return WriteInt16ToBitStream(int16(value64), ws)
			} else if value64 == int64(int32(value64)) {
				return WriteInt32ToBitStream(int32(value64), ws)
			} else {
				return WriteInt64ToBitStream(value64, ws)
			}
		} else {
			return errors.New("writeIntToBitStream unprefixed int error size detection")
		}
	default:
		return errors.New("writeAutoIntToBitStream only accepts int types (int8/16/32/64 and int)")
	}
}

// Manual Write No Mutations
func WriteInt8ToBitStream(value int8, ws *writeOnlyBitStream.WriteOnlyBitStream) error {
	if ws == nil {
		return errors.New("ws aka writeOnlyBitStream can't be nil")
	}
	//opCode for safe clone
	bitCodeToWrite := ubonOpCodes.UBON_OP_NEXT_INT.Clone()
	//mask 0 - 8 bits
	bitCodeToWrite.Append(false)
	//write single int8 byte
	bitCodeToWrite.AppendBitCode(bitcode.NewBitCodeFromBytes(byte(value)))
	//write to stream
	ws.AppendBitCode(bitCodeToWrite)
	return nil
}

func WriteInt16ToBitStream(value int16, ws *writeOnlyBitStream.WriteOnlyBitStream) error {
	if ws == nil {
		return errors.New("ws aka writeOnlyBitStream can't be nil")
	}
	//opCode for safe clone
	bitCodeToWrite := ubonOpCodes.UBON_OP_NEXT_INT.Clone()
	//mask 10 - 16 bits
	bitCodeToWrite.AppendBitCode(bitcode.NewBitCodeWithBoolCode(true, false))
	//write 2 bytes little endian
	bitCodeToWrite.AppendBitCode(bitcode.NewBitCodeFromBytes(byte(value), byte(value>>8)))
	//write to stream
	ws.AppendBitCode(bitCodeToWrite)
	return nil
}

func WriteInt32ToBitStream(value int32, ws *writeOnlyBitStream.WriteOnlyBitStream) error {
	if ws == nil {
		return errors.New("ws aka writeOnlyBitStream can't be nil")
	}
	//opCode for safe clone
	bitCodeToWrite := ubonOpCodes.UBON_OP_NEXT_INT.Clone()
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

func WriteInt64ToBitStream(value int64, ws *writeOnlyBitStream.WriteOnlyBitStream) error {
	if ws == nil {
		return errors.New("ws aka writeOnlyBitStream can't be nil")
	}
	//opCode for safe clone
	bitCodeToWrite := ubonOpCodes.UBON_OP_NEXT_INT.Clone()
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
