package ubonHeader

import (
	"errors"
	bitcode "ubon/internal/bitCode"
	"ubon/internal/readOnlyBitStream"
)

var TargetUBONSpecification_Bitcode = bitcode.NewBitCodeWithBoolCode(false)

type UbonHeader struct {
	Specification           bitcode.BitCode //1 bit now is 0 (prefixed 10 - spec 2 110 spec 3 etc)
	AlphabetSectionIncluded bool            //1 bit
}

func NewDefaultUbonHeader(isAlphabetIncluded bool) UbonHeader {
	return UbonHeader{
		Specification:           TargetUBONSpecification_Bitcode,
		AlphabetSectionIncluded: isAlphabetIncluded,
	}
}

func ReadUbonHeaderFromReadOnlyBitStream(bs *readOnlyBitStream.ReadOnlyBitStream) (*UbonHeader, error) {

	retVal := &UbonHeader{}
	specBitcode := bitcode.NewZeroBitCodeWithLength(0)
	for {
		b, err := bs.ReadBitCode(1)
		specBitcode.Append(b.GetBit(0))
		if err != nil {
			return nil, err
		}
		//zero end reading spec
		if !b.GetBit(0) {
			break
		}
	}
	//now single spec only
	if !specBitcode.Equal(TargetUBONSpecification_Bitcode) {
		return nil, errors.New("ubon specification unsupported")
	}
	retVal.Specification = specBitcode
	alphabetSectionIncluded, err := bs.ReadBitCode(1)
	if err != nil {
		return nil, err
	}
	retVal.AlphabetSectionIncluded = alphabetSectionIncluded.GetBit(0)
	return retVal, nil
}

func (h UbonHeader) BitCode() bitcode.BitCode {
	retVal := h.Specification.CloneAppend(h.AlphabetSectionIncluded)
	return retVal
}

func (h UbonHeader) Equal(other UbonHeader) bool {
	return h.Specification.Equal(other.Specification) &&
		h.AlphabetSectionIncluded == other.AlphabetSectionIncluded
}
