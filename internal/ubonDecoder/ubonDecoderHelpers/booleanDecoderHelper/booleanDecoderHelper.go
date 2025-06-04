package booleanDecoderHelper

import "ubon/internal/readOnlyBitStream"

func ReadBool(rs *readOnlyBitStream.ReadOnlyBitStream) (any, error) {
	boolBit, err := rs.ReadBitCode(1)
	if err != nil {
		return nil, err
	}
	return boolBit.GetBit(0), nil
}
