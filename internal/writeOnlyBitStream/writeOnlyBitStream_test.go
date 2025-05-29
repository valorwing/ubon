package writeOnlyBitStream_test

import (
	"testing"
	bitcode "ubon/internal/bitCode"
	"ubon/internal/writeOnlyBitStream"
)

func TestZero(t *testing.T) {

	bs := writeOnlyBitStream.NewWriteOnlyBitStream()
	bs.AppendBitCode(bitcode.NewBitCodeWithBoolCode(true))
	bs.AppendBitCode(bitcode.NewBitCodeWithBoolCode(true, false, true))
	bs.AppendBitCode(bitcode.NewBitCodeWithBoolCode(false, true))
	bs.AppendBitCode(bitcode.NewBitCodeWithBoolCode(false, true))

	bsBytes := bs.Bytes()
	if len(bsBytes) != 1 {
		t.Fail()
	}
	if bsBytes[0] != 0b11010101 {
		t.Fail()
	}

}
