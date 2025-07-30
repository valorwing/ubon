package ubon

import (
	"ubon/internal/ubonDecoder"
	"ubon/internal/ubonEncoder"
)

func MarshalUBON(input any) ([]byte, error) {

	return ubonEncoder.Encode(input)
}

func UnmarshalUBON(data []byte) (any, error) {

	return ubonDecoder.Decode(data)
}
