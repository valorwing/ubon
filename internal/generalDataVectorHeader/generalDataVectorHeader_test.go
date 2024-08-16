package generalDataVectorHeader_test

import (
	"log"
	"testing"
	"ubon/internal/generalDataVectorHeader"
)

func TestVersion(t *testing.T) {

	testHeader := generalDataVectorHeader.GeneralDataVectorHeader{
		ProtocolVersion: 1,
	}

	outData := make([]byte, 2)
	err := testHeader.WriteHeader(&outData)
	if err != nil {
		t.Fail()
	}

	restoredHeader, err := generalDataVectorHeader.ParseHeader(&outData)
	if err != nil {
		t.Fail()
	}
	if restoredHeader.ProtocolVersion != testHeader.ProtocolVersion {
		t.Fail()
	}
}

func TestEqual(t *testing.T) {

	testHeader := generalDataVectorHeader.GeneralDataVectorHeader{
		ProtocolVersion: 15,    // 4
		RootObjectType:  3,     // 2
		MaxLenBitsCount: 63,    // 6
		ReservedBit0:    false, // 1
		ReservedBit1:    true,  // 1
		ReservedBit2:    false,
		ReservedBit3:    true,
	}

	data := make([]byte, 2)
	err := testHeader.WriteHeader(&data)
	log.Printf("%+v", data)
	if err != nil {
		t.Fail()
	}
	restoredHeader, err := generalDataVectorHeader.ParseHeader(&data)
	if err != nil {
		t.Fail()
	}
	if !restoredHeader.IsEqualTo(&testHeader) {
		t.Fail()
	}
}
