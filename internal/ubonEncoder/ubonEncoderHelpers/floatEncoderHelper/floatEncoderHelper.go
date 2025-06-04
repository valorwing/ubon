package floatEncoderHelper

import (
	"errors"
	bitcode "ubon/internal/bitCode"
	"ubon/internal/ubonEncoder/ubonEncoderHelpers/intEncoderHelper"
	uIntEncoderHelper "ubon/internal/ubonEncoder/ubonEncoderHelpers/uintEncoderHelper"
	"ubon/internal/ubonOpCodes"
	"ubon/internal/writeOnlyBitStream"
	"unsafe"
)

// Auto write float and truncate and mutate to int value if needed
func WriteAutoFloatToBitStream(value any, ws *writeOnlyBitStream.WriteOnlyBitStream) error {

	switch value.(type) {
	case float32:
		float32Value := value.(float32)
		if float32Value == float32(int32(float32Value)) {
			//is int check sign and auto write int
			int32Value := int32(float32Value)
			if int32Value < 0 {
				return intEncoderHelper.WriteAutoIntToBitStream(int32Value, ws)
			} else {
				return uIntEncoderHelper.WriteAutoUnsignedIntToBitStream(uint32(int32Value), ws)
			}
		} else {
			return WriteFloat32ToBitStream(float32Value, ws)
		}
	case float64:
		float64Value := value.(float64)
		if float64Value == float64(int64(float64Value)) {
			//is int check sign and auto write int
			int64Value := int64(float64Value)
			if int64Value < 0 {
				return intEncoderHelper.WriteAutoIntToBitStream(int64Value, ws)
			} else {
				return uIntEncoderHelper.WriteAutoUnsignedIntToBitStream(uint64(int64Value), ws)
			}
		} else {
			return WriteFloat64ToBitStream(float64Value, ws)
		}
	default:
		return errors.New("writeAutoFloatToBitStream only support float32/float64")
	}

}

// Manual write float32 without mutations
func WriteFloat32ToBitStream(value float32, ws *writeOnlyBitStream.WriteOnlyBitStream) error {
	if ws == nil {
		return errors.New("ws aka writeOnlyBitStream can't be nil")
	}
	//write op code and flag is64bits
	ws.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_FLOAT.CloneAppend(false))
	//write 32 bits float
	valueBits := *(*uint32)(unsafe.Pointer(&value))
	ws.AppendBitCode(bitcode.NewBitCodeFromBytes(
		byte(valueBits), byte(valueBits>>8),
		byte(valueBits>>16), byte(valueBits>>24)))
	return nil
}

// Manual write float64 without mutations
func WriteFloat64ToBitStream(value float64, ws *writeOnlyBitStream.WriteOnlyBitStream) error {
	if ws == nil {
		return errors.New("ws aka writeOnlyBitStream can't be nil")
	}
	//write op code and flag is64bits
	ws.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_FLOAT.CloneAppend(true))
	//write 64 bits float
	valueBits := *(*uint64)(unsafe.Pointer(&value))
	ws.AppendBitCode(bitcode.NewBitCodeFromBytes(
		byte(valueBits), byte(valueBits>>8),
		byte(valueBits>>16), byte(valueBits>>24),
		byte(valueBits>>32), byte(valueBits>>40),
		byte(valueBits>>48), byte(valueBits>>56)))
	return nil
}
