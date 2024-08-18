package ubonTypes_test

import (
	"runtime"
	"testing"
	"ubon/internal/ubonTypes.go"
)

func is64Bit() bool {

	var arch64 = []string{
		"amd64",
		"arm64",
		"arm64be",
		"ppc64",
		"ppc64le",
		"mips64",
		"mips64le",
		"riscv64",
		"s390x",
		"sparc64",
		"wasm"}

	is64BitVal := false
	for _, arch := range arch64 {
		if arch == runtime.GOARCH {
			is64BitVal = true
			break
		}
	}

	return is64BitVal
}
func TestDetectDefaultIntsLen(t *testing.T) {

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

func TestFloat64_28_TypeDetect(t *testing.T) {

	var testData = float64(100)
	detectedType, err := ubonTypes.DetectType(testData)
	if err != nil {
		t.Fail()
	}
	if detectedType != ubonTypes.UBON_Float64_40 {
		t.Fail()
	}

	testData = float64(505)
	detectedType, err = ubonTypes.DetectType(testData)
	if err != nil {
		t.Fail()
	}
	if detectedType != ubonTypes.UBON_Float64_40 {
		t.Fail()
	}

}

func TestUnions(t *testing.T) {

	a := 10
	b := map[string]interface{}{}
	aType, err := ubonTypes.DetectType(a)
	if err != nil {
		t.Fail()
	}
	if aType.GetUnion() != ubonTypes.Primitives {
		t.Fail()
	}
	bType, err := ubonTypes.DetectType(b)
	if err != nil {
		t.Fail()
	}
	if bType.GetUnion() != ubonTypes.Objects {
		t.Fail()
	}

}
