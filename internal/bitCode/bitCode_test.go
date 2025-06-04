package bitcode_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	bitcode "ubon/internal/bitCode"
)

func TestHash(t *testing.T) {
	aBytes := []byte{0, 134, 255}
	aBitcode := bitcode.NewBitCodeFromBytes(aBytes...)
	bBitcode := aBitcode.Clone()
	if aBitcode.Hash() != bBitcode.Hash() {
		t.Fail()
	}
	bBitcode.Append(false)
	if aBitcode.Hash() == bBitcode.Hash() {
		t.Fail()
	}

	//empty
	aBitcode = bitcode.NewZeroBitCodeWithLength(0)
	bBitcode = bitcode.NewZeroBitCodeWithLength(0)
	if aBitcode.Hash() != bBitcode.Hash() {
		t.Fail()
	}
}

func TestByteEquality(t *testing.T) {
	aBytes := []byte{0, 134, 255}
	aBitcode := bitcode.NewBitCodeFromBytes(aBytes...)
	if !bytes.Equal(aBytes, aBitcode.Bytes()) {
		t.Fail()
	}
}

func TestAppend(t *testing.T) {
	a := bitcode.NewBitCodeWithBoolCode(true, true, true, false)

	sTest := strings.Repeat(a.String(), 2)
	b := a.Clone()

	fmt.Println(a.String())
	fmt.Println(b.String())

	b.AppendBitCode(a)
	bStr := b.String()
	if bStr != sTest {
		t.Fail()
	}
}

func TestFromBytes(t *testing.T) {
	a := []byte{255, 255}
	b := []byte{0, 0, 0, 0, 255}

	aStr := strings.Repeat("1", 16)
	bStr := strings.Repeat("0", 32) + strings.Repeat("1", 8)

	aBitCode := bitcode.NewBitCodeFromBytes(a...)
	bBitCode := bitcode.NewBitCodeFromBytes(b...)

	if aStr != aBitCode.String() {
		t.Fail()
	}
	if bStr != bBitCode.String() {
		t.Fail()
	}
	aBitCode.Append(false)
	aStr += "0"
	aBitCodeStr := aBitCode.String()
	if aBitCodeStr != aStr {
		t.Fail()
	}
}

func TestZero(t *testing.T) {

	a := bitcode.NewBitCodeWithBoolCode(true, true, true, false)
	aStr := a.String()
	if aStr != "1110" {
		t.Fail()
	}

	if a.BitLength() != 4 {
		t.Fail()
	}

	a.Append(true)

	aStr = a.String()
	if aStr != "11101" {
		t.Fail()
	}
	if a.BitLength() != 5 {
		t.Fail()
	}

	a.Append(false)

	aStr = a.String()
	if aStr != "111010" {
		t.Fail()
	}
	if a.BitLength() != 6 {
		t.Fail()
	}
}
