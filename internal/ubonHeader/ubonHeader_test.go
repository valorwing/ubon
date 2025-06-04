package ubonHeader_test

import (
	"testing"
	"ubon/internal/readOnlyBitStream"
	"ubon/internal/ubonHeader"
	"ubon/internal/writeOnlyBitStream"
)

func TestZero(t *testing.T) {
	header := ubonHeader.UbonHeader{
		Specification:           ubonHeader.TargetUBONSpecification_Bitcode,
		AlphabetSectionIncluded: false,
	}
	ws := writeOnlyBitStream.NewWriteOnlyBitStream()
	ws.AppendBitCode(header.BitCode())
	rs := readOnlyBitStream.NewReadOnlyBitStream(ws.Bytes())
	restoredHeader, err := ubonHeader.ReadUbonHeaderFromReadOnlyBitStream(&rs)
	if err != nil {
		t.Fail()
		return
	}
	if !restoredHeader.Equal(header) {
		t.Fail()
	}
}
