package uIntDecoderHelper

import (
	"errors"
	bitcode "ubon/internal/bitCode"
	"ubon/internal/readOnlyBitStream"
)

func ReadMaskedAutoUnsignedInt(rs *readOnlyBitStream.ReadOnlyBitStream) (any, error) {
	bitMask := bitcode.NewZeroBitCodeWithLength(0)
	for {
		b, err := rs.ReadBitCode(1)
		if err != nil {
			return nil, err
		}
		bitMask.AppendBitCode(*b)
		//read first zero - end or max mask len
		if !b.GetBit(0) || bitMask.BitLength() == 3 {
			break
		}
	}
	//0 - 8 bit
	if bitMask.Equal(bitcode.NewBitCodeWithBoolCode(false)) {
		return ReadUInt8(rs)
		//01 - 16 bit
	} else if bitMask.Equal(bitcode.NewBitCodeWithBoolCode(true, false)) {
		return ReadUInt16(rs)
	} else if bitMask.Equal(bitcode.NewBitCodeWithBoolCode(true, true, false)) {
		return ReadUInt32(rs)
	} else if bitMask.Equal(bitcode.NewBitCodeWithBoolCode(true, true, true)) {
		return ReadUInt64(rs)
	} else {
		return nil, errors.New("readMaskedAutoInt invalid bit-length mask")
	}
}

func ReadUInt8(rs *readOnlyBitStream.ReadOnlyBitStream) (any, error) {
	int8BitCode, err := rs.ReadBitCode(8)
	if err != nil {
		return nil, err
	}
	return uint8(int8BitCode.Bytes()[0]), nil
}

func ReadUInt16(rs *readOnlyBitStream.ReadOnlyBitStream) (any, error) {
	int16BitCode, err := rs.ReadBitCode(16)
	if err != nil {
		return nil, err
	}
	int16Bytes := int16BitCode.Bytes()
	return uint16(int16Bytes[0]) | uint16(int16Bytes[1])<<8, nil
}

func ReadUInt32(rs *readOnlyBitStream.ReadOnlyBitStream) (any, error) {
	int32BitCode, err := rs.ReadBitCode(32)
	if err != nil {
		return nil, err
	}
	int32Bytes := int32BitCode.Bytes()
	return uint32(int32Bytes[0]) | uint32(int32Bytes[1])<<8 | uint32(int32Bytes[2])<<16 | uint32(int32Bytes[3])<<24, nil
}

func ReadUInt64(rs *readOnlyBitStream.ReadOnlyBitStream) (any, error) {
	int64BitCode, err := rs.ReadBitCode(64)
	if err != nil {
		return nil, err
	}
	int64Bytes := int64BitCode.Bytes()

	return uint64(int64Bytes[0]) | uint64(int64Bytes[1])<<8 | uint64(int64Bytes[2])<<16 | uint64(int64Bytes[3])<<24 |
		uint64(int64Bytes[4])<<32 | uint64(int64Bytes[5])<<40 | uint64(int64Bytes[6])<<48 | uint64(int64Bytes[7])<<56, nil
}
