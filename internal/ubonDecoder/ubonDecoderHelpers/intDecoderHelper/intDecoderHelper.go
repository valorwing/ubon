package intDecoderHelper

import (
	"errors"
	bitcode "ubon/internal/bitCode"
	"ubon/internal/readOnlyBitStream"
)

func ReadMaskedAutoInt(rs *readOnlyBitStream.ReadOnlyBitStream) (any, error) {
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
		return ReadInt8(rs)
		//01 - 16 bit
	} else if bitMask.Equal(bitcode.NewBitCodeWithBoolCode(true, false)) {
		return ReadInt16(rs)
	} else if bitMask.Equal(bitcode.NewBitCodeWithBoolCode(true, true, false)) {
		return ReadInt32(rs)
	} else if bitMask.Equal(bitcode.NewBitCodeWithBoolCode(true, true, true)) {
		return ReadInt64(rs)
	} else {
		return nil, errors.New("readMaskedAutoInt invalid bit-length mask")
	}
}

func ReadInt8(rs *readOnlyBitStream.ReadOnlyBitStream) (any, error) {
	int8BitCode, err := rs.ReadBitCode(8)
	if err != nil {
		return nil, err
	}
	return int8(int8BitCode.Bytes()[0]), nil
}

func ReadInt16(rs *readOnlyBitStream.ReadOnlyBitStream) (any, error) {
	int16BitCode, err := rs.ReadBitCode(16)
	if err != nil {
		return nil, err
	}
	int16Bytes := int16BitCode.Bytes()
	return int16(int16Bytes[0]) | int16(int16Bytes[1])<<8, nil
}

func ReadInt32(rs *readOnlyBitStream.ReadOnlyBitStream) (any, error) {
	int32BitCode, err := rs.ReadBitCode(32)
	if err != nil {
		return nil, err
	}
	int32Bytes := int32BitCode.Bytes()
	return int32(int32Bytes[0]) | int32(int32Bytes[1])<<8 | int32(int32Bytes[2])<<16 | int32(int32Bytes[3])<<24, nil
}

func ReadInt64(rs *readOnlyBitStream.ReadOnlyBitStream) (any, error) {
	int64BitCode, err := rs.ReadBitCode(64)
	if err != nil {
		return nil, err
	}
	int64Bytes := int64BitCode.Bytes()

	return int64(int64Bytes[0]) | int64(int64Bytes[1])<<8 | int64(int64Bytes[2])<<16 | int64(int64Bytes[3])<<24 |
		int64(int64Bytes[4])<<32 | int64(int64Bytes[5])<<40 | int64(int64Bytes[6])<<48 | int64(int64Bytes[7])<<56, nil
}
