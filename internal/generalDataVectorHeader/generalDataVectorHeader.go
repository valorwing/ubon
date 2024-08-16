package generalDataVectorHeader

import (
	"fmt"
	"ubon/internal/bitutil"
	"ubon/internal/errors"
)

type GeneralDataVectorHeader struct { // 16 bit | 2 byte
	ProtocolVersion uint8 // 4 bit | 0 - 15
	RootObjectType  uint8 // 0-3 | 2 bit | 1 - primitive, 2 - array, 3 - object
	MaxLenBitsCount uint8 // 6 bit | 0-63 -> 1 - 64
	ReservedBit0    bool  // default must be false
	ReservedBit1    bool  // default must be false
	ReservedBit2    bool  // default must be false
	ReservedBit3    bool  // default must be false
}

func (h *GeneralDataVectorHeader) IsEqualTo(h1 *GeneralDataVectorHeader) bool {
	return h.ProtocolVersion == h1.ProtocolVersion &&
		h.RootObjectType == h1.RootObjectType &&
		h.MaxLenBitsCount == h1.MaxLenBitsCount &&
		h.ReservedBit0 == h1.ReservedBit0 &&
		h.ReservedBit1 == h1.ReservedBit1 &&
		h.ReservedBit2 == h1.ReservedBit2 &&
		h.ReservedBit3 == h1.ReservedBit3
}
func (h *GeneralDataVectorHeader) WriteHeader(output *[]byte) error {
	var err error
	out := *output
	if len(out) < 2 {
		err = errors.InvalidOutDataPointer("out len for header must be > 2. current len is:" + fmt.Sprintf("%d", len(out)))
		return err
	}

	out[0] = (h.ProtocolVersion << 4) | ((h.RootObjectType & 0b00000011) << 2)
	out[0] |= ((h.MaxLenBitsCount >> 4) & 0b00000011)

	out[1] = (h.MaxLenBitsCount << 4)

	if h.ReservedBit0 {
		out[1] = bitutil.SetBit(out[1], 4)
	} else {
		out[1] = bitutil.ResetBit(out[1], 4)
	}
	if h.ReservedBit1 {
		out[1] = bitutil.SetBit(out[1], 5)
	} else {
		out[1] = bitutil.ResetBit(out[1], 5)
	}
	if h.ReservedBit2 {
		out[1] = bitutil.SetBit(out[1], 6)
	} else {
		out[1] = bitutil.ResetBit(out[1], 6)
	}
	if h.ReservedBit3 {
		out[1] = bitutil.SetBit(out[1], 7)
	} else {
		out[1] = bitutil.ResetBit(out[1], 7)
	}

	return err
}

func ParseHeader(input *[]byte) (*GeneralDataVectorHeader, error) {
	var err error
	retVal := &GeneralDataVectorHeader{}
	in := *input
	if len(in) < 2 {
		err = errors.InvalidHeader("current len is: " + fmt.Sprintf("%d", len(in)) + " minimum allowed len header 2 byte")
		return retVal, err
	}

	retVal.ProtocolVersion = in[0] >> 4
	retVal.RootObjectType = (in[0] & 0b00001100) >> 2
	retVal.MaxLenBitsCount = ((in[0] & 0b00000011) << 4) | (in[1] >> 4)

	retVal.ReservedBit0 = bitutil.ReadBit(in[1], 4)
	retVal.ReservedBit1 = bitutil.ReadBit(in[1], 5)
	retVal.ReservedBit2 = bitutil.ReadBit(in[1], 6)
	retVal.ReservedBit3 = bitutil.ReadBit(in[1], 7)

	return retVal, err
}
