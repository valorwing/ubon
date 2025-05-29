package bitcode

import "testing"

func TestZero(t *testing.T) {

	a := NewBitCodeWithBoolCode(true, true, true, false)
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
