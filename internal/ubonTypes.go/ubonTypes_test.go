package ubonTypes_test

import (
	"log"
	"runtime"
	"strings"
	"testing"
	"ubon/internal/ubonTypes.go"
)

func is64Bit() bool {
	return runtime.GOARCH == "amd64" || runtime.GOARCH == "arm64"
}
func TestDetectDefaultIntsLen(t *testing.T) {

	log.Println(strings.Repeat("0", 64-28) + strings.Repeat("1", 28))

	if is64Bit() {
		if !ubonTypes.GetDefaultIntIs32BitsLen() && !ubonTypes.GetDefaultUIntIs32BitsLen() {
			t.Log("Host 64 bit. Len not 32 is ok")
		} else {
			t.Fail()
		}
	} else {
		if ubonTypes.GetDefaultIntIs32BitsLen() && ubonTypes.GetDefaultUIntIs32BitsLen() {
			t.Log("Host 64 bit. Len is 32 is ok")
		} else {
			t.Fail()
		}
	}

}

func TestDetectTypeBoolBase(t *testing.T) {

	var boolVar interface{} = false
	detectedType, err := ubonTypes.DetectType(boolVar)
	if err != nil {
		t.Fail()
	}
	if detectedType != ubonTypes.UBON_Bool {
		t.Fail()
	}
}

func TestIntPointer(t *testing.T) {

	var a = uint8(255)
	var aPtr *uint8 = &a

	ubonType, err := ubonTypes.DetectType(aPtr)
	if err != nil {
		t.Fail()
	}
	if ubonType != ubonTypes.UBON_UInt8 {
		t.Fail()
	}
}

func TestInt64Compression(t *testing.T) {

	var testData = uint(1)
	detectedType, err := ubonTypes.DetectType(testData)

	if err != nil {
		t.Fail()
	}
	if detectedType != ubonTypes.UBON_UInt4 {
		t.Fail()
	}

	testData1 := int(0xFFFF)

	detectedType, err = ubonTypes.DetectType(testData1)
	if err != nil {
		t.Fail()
	}
	if detectedType != ubonTypes.UBON_Int16 {
		t.Fail()
	}

	detectedType, err = ubonTypes.DetectType(testData1 + 1)
	if err != nil {
		t.Fail()
	}
	if detectedType != ubonTypes.UBON_Int24 {
		t.Fail()
	}
}

func TestUIntCompression(t *testing.T) {

	var testData4 = uint8(15)
	detectedType, err := ubonTypes.DetectType(testData4)
	if err != nil {
		t.Fail()
	}
	if detectedType != ubonTypes.UBON_UInt4 {
		t.Fail()
	}
	var testData8 interface{} = uint8(16)
	detectedType, err = ubonTypes.DetectType(testData8)
	if err != nil {
		t.Fail()
	}
	if detectedType != ubonTypes.UBON_UInt8 {
		t.Fail()
	}
}
