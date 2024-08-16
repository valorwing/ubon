package bitutil_test

import (
	"bytes"
	"testing"
	"ubon/internal/bitutil"
)

func TestBase(t *testing.T) {

	testData := []byte{0, 0, 128, 255}

	testBuffer := [16]byte{}
	copy(testData, testBuffer[:])

	testOutput := make([]byte, 4)

	bitutil.WriteBitsPrimitive(&testBuffer, 32, &testOutput, 0, 0)

	if !bytes.Equal(testData, testOutput) {
		t.Fail()
	}

}

func TestAppendWrite(t *testing.T) {

	testDataA := []byte{10}
	testDataB := []byte{20}
	testBuffer := [16]byte{}
	testOutput := make([]byte, 2)
	copy(testDataA, testBuffer[:])
	bitutil.WriteBitsPrimitive(&testBuffer, 8, &testOutput, 0, 0)
	copy(testDataB, testBuffer[:])
	bitutil.WriteBitsPrimitive(&testBuffer, 8, &testOutput, 1, 0)
	if !bytes.Equal(append(testDataA, testDataB...), testOutput) {
		t.Fail()
	}

}

func TestAppendWrite4NumbersInSingleByte(t *testing.T) {
	testDataA := []byte{1}
	testDataB := []byte{2}
	testDataC := []byte{3}
	testDataD := []byte{0}
	testOutput := make([]byte, 1)
	testBuffer := [16]byte{}

	copy(testDataA, testBuffer[:])
	bitutil.WriteBitsPrimitive(&testBuffer, 2, &testOutput, 0, 0)

	copy(testDataB, testBuffer[:])
	bitutil.WriteBitsPrimitive(&testBuffer, 2, &testOutput, 0, 2)

	copy(testDataC, testBuffer[:])
	bitutil.WriteBitsPrimitive(&testBuffer, 2, &testOutput, 0, 4)

	copy(testDataD, testBuffer[:])
	bitutil.WriteBitsPrimitive(&testBuffer, 2, &testOutput, 0, 6)

	restoreDataA := make([]byte, 1)
	restoreDataB := make([]byte, 1)
	restoreDataC := make([]byte, 1)
	restoreDataD := make([]byte, 1)
	bitutil.ReadBitsPrimitive(&testBuffer, 2, &testOutput, 0, 0)
	copy(restoreDataA, testBuffer[:1])
	bitutil.ReadBitsPrimitive(&testBuffer, 2, &testOutput, 0, 2)
	copy(restoreDataB, testBuffer[:1])
	bitutil.ReadBitsPrimitive(&testBuffer, 2, &testOutput, 0, 4)
	copy(restoreDataC, testBuffer[:1])
	bitutil.ReadBitsPrimitive(&testBuffer, 2, &testOutput, 0, 6)
	copy(restoreDataD, testBuffer[:1])

	if !bytes.Equal(testDataA, restoreDataA) {
		t.Fail()
	}
	if !bytes.Equal(testDataB, restoreDataB) {
		t.Fail()
	}
	if !bytes.Equal(testDataC, restoreDataC) {
		t.Fail()
	}
	if !bytes.Equal(testDataD, restoreDataD) {
		t.Fail()
	}
}
