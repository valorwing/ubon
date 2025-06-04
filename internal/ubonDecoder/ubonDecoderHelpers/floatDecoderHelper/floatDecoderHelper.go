package floatDecoderHelper

import (
	"ubon/internal/readOnlyBitStream"
	"unsafe"
)

func ReadMaskedAutoFloat(rs *readOnlyBitStream.ReadOnlyBitStream) (any, error) {
	//Mask 0 - float32, 1 - float64
	maskBitCode, err := rs.ReadBitCode(1)
	if err != nil {
		return nil, err
	}
	if !maskBitCode.GetBit(0) {
		return ReadFloat32(rs)
	} else {
		return ReadFloat64(rs)
	}
}

func ReadFloat32(rs *readOnlyBitStream.ReadOnlyBitStream) (any, error) {
	float32Bitcode, err := rs.ReadBitCode(32)
	if err != nil {
		return nil, err
	}
	float32Bytes := float32Bitcode.Bytes()
	float32Bits := uint32(float32Bytes[0]) | uint32(float32Bytes[1])<<8 | uint32(float32Bytes[2])<<16 | uint32(float32Bytes[3])<<24
	return *(*float32)(unsafe.Pointer(&float32Bits)), nil
}

func ReadFloat64(rs *readOnlyBitStream.ReadOnlyBitStream) (any, error) {
	float64Bitcode, err := rs.ReadBitCode(64)
	if err != nil {
		return nil, err
	}
	float64Bytes := float64Bitcode.Bytes()
	float64Bits := uint64(float64Bytes[0]) | uint64(float64Bytes[1])<<8 | uint64(float64Bytes[2])<<16 | uint64(float64Bytes[3])<<24 |
		uint64(float64Bytes[4])<<32 | uint64(float64Bytes[5])<<40 | uint64(float64Bytes[6])<<48 | uint64(float64Bytes[7])<<56
	return *(*float64)(unsafe.Pointer(&float64Bits)), nil
}
