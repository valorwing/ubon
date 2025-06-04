package ubonEncoder_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
	"ubon/internal/ubonEncoder"

	"github.com/fxamacker/cbor/v2"
)

func TestEncodeSingleNull(t *testing.T) {

	var obj interface{} = nil
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		t.Fail()
	}
	ubonBytes, err := ubonEncoder.Encode(obj)
	if err != nil {
		t.Fail()
	}
	fmt.Println("UBON bytes len: ", len(ubonBytes))
	fmt.Println("JSON bytes len: ", len(jsonBytes))
}

func TestEncodeSingleBoolean(t *testing.T) {
	var testData bool = true
	jsonBytes, err := json.Marshal(testData)
	if err != nil {
		t.Fail()
	}
	ubonBytes, err := ubonEncoder.Encode(testData)
	if err != nil {
		t.Fail()
	}
	fmt.Println("UBON bytes len: ", len(ubonBytes))
	fmt.Println("JSON bytes len: ", len(jsonBytes))
}

func TestEncodeSingleShortString(t *testing.T) {
	var testData string = "hello world"
	jsonBytes, err := json.Marshal(testData)
	if err != nil {
		t.Fail()
	}
	ubonBytes, err := ubonEncoder.Encode(testData)
	if err != nil {
		t.Fail()
	}
	fmt.Println("UBON bytes len: ", len(ubonBytes))
	fmt.Println("JSON bytes len: ", len(jsonBytes))
}

func TestEncodeSingleLongString(t *testing.T) {
	var testData string = `Once upon a time, a small fox named Felix found a forgotten flute in the forest. Curious, he played a note. To his surprise, the trees shimmered, the wind paused, and birds began to sing in harmony. "Magic!" he whispered. Every day after, the forest danced to Felix's tune.`
	jsonBytes, err := json.Marshal(testData)
	if err != nil {
		t.Fail()
	}
	ubonBytes, err := ubonEncoder.Encode(testData)
	if err != nil {
		t.Fail()
	}
	fmt.Println("UBON bytes len: ", len(ubonBytes))
	fmt.Println("JSON bytes len: ", len(jsonBytes))
}

func TestEncodeSingleInt8(t *testing.T) {
	var testData int64 = 6
	jsonBytes, err := json.Marshal(testData)
	if err != nil {
		t.Fail()
	}
	ubonBytes, err := ubonEncoder.Encode(testData)
	if err != nil {
		t.Fail()
	}
	fmt.Println("UBON bytes len: ", len(ubonBytes))
	fmt.Println("JSON bytes len: ", len(jsonBytes))
}

func TestEncodeSingleInt64(t *testing.T) {
	var testData int64 = 0x7FFFFFFFFFFFFFFF
	jsonBytes, err := json.Marshal(testData)
	if err != nil {
		t.Fail()
	}
	ubonBytes, err := ubonEncoder.Encode(testData)
	if err != nil {
		t.Fail()
	}
	fmt.Println("UBON bytes len: ", len(ubonBytes))
	fmt.Println("JSON bytes len: ", len(jsonBytes))
}

func TestEncodeSingleUInt8(t *testing.T) {
	var testData uint8 = 50
	jsonBytes, err := json.Marshal(testData)
	if err != nil {
		t.Fail()
	}
	ubonBytes, err := ubonEncoder.Encode(testData)
	if err != nil {
		t.Fail()
	}
	fmt.Println("UBON bytes len: ", len(ubonBytes))
	fmt.Println("JSON bytes len: ", len(jsonBytes))
}

func TestEncodeSingleUInt64(t *testing.T) {
	var testData uint64 = 0x7FFFFFFFFFFFFFFF
	jsonBytes, err := json.Marshal(testData)
	if err != nil {
		t.Fail()
	}
	ubonBytes, err := ubonEncoder.Encode(testData)
	if err != nil {
		t.Fail()
	}
	fmt.Println("UBON bytes len: ", len(ubonBytes))
	fmt.Println("JSON bytes len: ", len(jsonBytes))
}

func TestEncodeSingleFloat64(t *testing.T) {
	var testData float64 = 3.141592653589793115997963468544185161590576171875
	jsonBytes, err := json.Marshal(testData)
	if err != nil {
		t.Fail()
	}
	ubonBytes, err := ubonEncoder.Encode(testData)
	if err != nil {
		t.Fail()
	}
	fmt.Println("UBON bytes len: ", len(ubonBytes))
	fmt.Println("JSON bytes len: ", len(jsonBytes))
}

func TestEncodeSingleFloat32(t *testing.T) {
	var testData float32 = 3.141592653589

	jsonBytes, err := json.Marshal(testData)
	if err != nil {
		t.Fail()
	}
	ubonBytes, err := ubonEncoder.Encode(testData)
	if err != nil {
		t.Fail()
	}
	fmt.Println("UBON bytes len: ", len(ubonBytes))
	fmt.Println("JSON bytes len: ", len(jsonBytes))
}

func TestEncodeJSON(t *testing.T) {
	rawJSON := `
	{
  "nullValue": null,
  "booleanTrue": true,
  "booleanFalse": false,
  "integer": 123,
  "negativeInteger": -456,
  "float": 78.9,
  "negativeFloat": -0.001,
  "zero": 0,
  "emptyString": "",
  "string": "Hello, 世界!",
  "escapedString": "Line1\nLine2\tTabbed\\Backslash\"Quote",
  "nestedObject": {
    "level1Key": "level1Value",
    "level1Null": null,
    "level2Object": {
      "level2Key": 42,
      "level2Null": null,
      "level3Object": {
        "deepKey": false,
        "deepNull": null,
        "anotherString": "deep"
      }
    },
    "anotherNull": null
  },
  "trueString": "true",
  "falseString": "false",
  "nullString": "null",
  "numberAsString": "12345"
	}
	`

	var jsonSerializedObject map[string]any
	var jsonSerializedObjectCopy map[string]any

	err := json.Unmarshal([]byte(rawJSON), &jsonSerializedObject)
	if err != nil {
		t.Fail()
	}

	err = json.Unmarshal([]byte(rawJSON), &jsonSerializedObjectCopy)
	if err != nil {
		t.Fail()
	}

	timeStartUBON := time.Now()
	ubonBytes, err := ubonEncoder.Encode(jsonSerializedObject)
	ubonDuration := time.Since(timeStartUBON)
	if err != nil {
		t.FailNow()
		return
	}
	timeStartJSON := time.Now()
	compactJsonBytes, err := json.Marshal(jsonSerializedObjectCopy)
	jsonDuration := time.Since(timeStartJSON)
	if err != nil {
		t.FailNow()
		return
	}

	timeStartCBOR := time.Now()
	cborBytes, err := cbor.Marshal(jsonSerializedObjectCopy)
	cborDuration := time.Since(timeStartCBOR)
	if err != nil {
		t.FailNow()
		return
	}

	fmt.Println("ubon bytes len: ", len(ubonBytes))
	fmt.Println("ubon serialization duratuion: ", ubonDuration)
	fmt.Println("cbor bytes len: ", len(cborBytes))
	fmt.Println("cbor serialization duratuion: ", cborDuration)
	fmt.Println("compact json bytes len: ", len(compactJsonBytes))
	fmt.Println("json serialization duration: ", jsonDuration)
}
