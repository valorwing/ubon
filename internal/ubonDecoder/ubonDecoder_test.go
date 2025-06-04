package ubonDecoder_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"testing"
	"time"
	"ubon/internal/ubonDecoder"
	"ubon/internal/ubonEncoder"

	"github.com/fxamacker/cbor/v2"
)

func TestNullAndCorrectNullEquality(t *testing.T) {
	var obj interface{} = nil
	ubonBytes, err := ubonEncoder.Encode(obj)
	if err != nil {
		t.Fail()
	}
	fmt.Println("ubon len: ", len(ubonBytes))
	decodedObj, err := ubonDecoder.Decode(ubonBytes)
	if err != nil {
		t.Fail()
	}
	if decodedObj != obj {
		t.Fail()
	}
	if decodedObj != nil {
		t.Fail()
	}
}

func TestSingleBoolEqualityPositive(t *testing.T) {
	var obj bool = true
	ubonBytes, err := ubonEncoder.Encode(obj)
	if err != nil {
		t.Fail()
	}
	fmt.Println("ubon len: ", len(ubonBytes))
	decodedObj, err := ubonDecoder.Decode(ubonBytes)
	if err != nil {
		t.Fail()
	}
	if decodedObj != obj {
		t.Fail()
	}
}

func TestSingleBoolEqualityNegative(t *testing.T) {
	var obj bool = false
	ubonBytes, err := ubonEncoder.Encode(obj)
	if err != nil {
		t.Fail()
	}
	fmt.Println("ubon len: ", len(ubonBytes))
	decodedObj, err := ubonDecoder.Decode(ubonBytes)
	if err != nil {
		t.Fail()
	}
	if decodedObj != obj {
		t.Fail()
	}
}

func TestRoundTrip(t *testing.T) {
	cases := []any{
		nil,
		true, false,
		42, int64(1 << 40),
		uint(123456),
		3.14, float64(2.718),
		"hello world",
	}

	for _, input := range cases {
		encoded, err := ubonEncoder.Encode(input)
		if err != nil {
			t.Errorf("encode failed: %v", err)
			continue
		}
		decoded, err := ubonDecoder.Decode(encoded)
		if err != nil {
			t.Errorf("decode failed: %v", err)
			continue
		}
		if fmt.Sprintf("%v", decoded) != fmt.Sprintf("%v", input) {
			t.Errorf("roundtrip mismatch: input=%v decoded=%v", input, decoded)
		}
	}
}

func TestSingleStringEquality(t *testing.T) {
	var obj string = `Once upon a time, a small fox named Felix found a forgotten flute in the forest. Curious, he played a note. To his surprise, the trees shimmered, the wind paused, and birds began to sing in harmony. "Magic!" he whispered. Every day after, the forest danced to Felix's tune.
`
	ubonBytes, err := ubonEncoder.Encode(obj)
	if err != nil {
		t.Fail()
	}
	fmt.Println("ubon len: ", len(ubonBytes))

	decodedObj, err := ubonDecoder.Decode(ubonBytes)
	if err != nil {
		t.Fail()
	}
	if decodedObj != obj {
		t.Fail()
	}
}

func TestSingleInt8(t *testing.T) {
	//test auto transfrom
	var obj int64 = 100
	ubonBytes, err := ubonEncoder.Encode(obj)
	if err != nil {
		t.Fail()
	}
	fmt.Println("ubon len: ", len(ubonBytes))
	decodedObj, err := ubonDecoder.Decode(ubonBytes)
	if err != nil {
		t.Fail()
	}

	if obj != int64(decodedObj.(uint8)) {
		t.Fail()
	}
}

func TestSingleInt16(t *testing.T) {
	//test auto transfrom
	var obj int64 = 512
	ubonBytes, err := ubonEncoder.Encode(obj)
	if err != nil {
		t.Fail()
	}
	fmt.Println("ubon len: ", len(ubonBytes))
	decodedObj, err := ubonDecoder.Decode(ubonBytes)
	if err != nil {
		t.Fail()
	}

	if obj != int64(decodedObj.(uint16)) {
		t.Fail()
	}
}

func TestSingleInt64(t *testing.T) {
	//test auto transfrom
	var obj int64 = 0x7FFFFFFFFFFFFFFF
	ubonBytes, err := ubonEncoder.Encode(obj)
	if err != nil {
		t.Fail()
	}
	fmt.Println("ubon len: ", len(ubonBytes))
	decodedObj, err := ubonDecoder.Decode(ubonBytes)
	if err != nil {
		t.Fail()
	}

	if obj != int64(decodedObj.(uint64)) {
		t.Fail()
	}
}

func TestSingleUInt16(t *testing.T) {
	//test auto transfrom
	var obj uint64 = 512
	ubonBytes, err := ubonEncoder.Encode(obj)
	if err != nil {
		t.Fail()
	}
	fmt.Println("ubon len: ", len(ubonBytes))
	decodedObj, err := ubonDecoder.Decode(ubonBytes)
	if err != nil {
		t.Fail()
	}

	if obj != uint64(decodedObj.(uint16)) {
		t.Fail()
	}
}

func TestSingleUInt64(t *testing.T) {
	//test auto transfrom
	var obj uint64 = 0x7FFFFFFFFFFFFFFF
	ubonBytes, err := ubonEncoder.Encode(obj)
	if err != nil {
		t.Fail()
	}
	fmt.Println("ubon len: ", len(ubonBytes))
	decodedObj, err := ubonDecoder.Decode(ubonBytes)
	if err != nil {
		t.Fail()
	}

	if obj != decodedObj {
		t.Fail()
	}
}

func TestSingleFloat64(t *testing.T) {
	//test auto transfrom
	var obj float64 = 3.141592653589793115997963468544185161590576171875
	ubonBytes, err := ubonEncoder.Encode(obj)
	if err != nil {
		t.Fail()
	}
	fmt.Println("ubon len: ", len(ubonBytes))
	decodedObj, err := ubonDecoder.Decode(ubonBytes)
	if err != nil {
		t.Fail()
	}

	if obj != decodedObj {
		t.Fail()
	}
}

func TestSingleFloat32(t *testing.T) {
	//test auto transfrom
	var obj float32 = 3.1415926535
	ubonBytes, err := ubonEncoder.Encode(obj)
	if err != nil {
		t.Fail()
	}
	fmt.Println("ubon len: ", len(ubonBytes))
	decodedObj, err := ubonDecoder.Decode(ubonBytes)
	if err != nil {
		t.Fail()
	}

	if obj != decodedObj {
		t.Fail()
	}
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}

func generateBigJSON() map[string]interface{} {
	rand.Seed(42)
	data := make(map[string]interface{})

	// Добавляем 50 вложенных уровней
	current := data
	for i := 0; i < 50; i++ {
		nested := make(map[string]interface{})
		current[fmt.Sprintf("level%d", i)] = nested

		// Заполняем каждый уровень
		for j := 0; j < 500; j++ {
			key := fmt.Sprintf("key%04d", j)

			switch rand.Intn(6) {
			case 0:
				nested[key] = nil
			case 1:
				nested[key] = rand.Float64() * 1000000
			case 2:
				nested[key] = rand.Intn(1000000) - 500000
			case 3:
				nested[key] = String(500) // Случайные строки
			case 4:
				nested[key] = rand.Intn(2) == 1
			case 5:
				// Вложенные бинарные данные как Base64
				buf := make([]byte, 128)
				rand.Read(buf)
				nested[key] = base64.StdEncoding.EncodeToString(buf)
			}
		}
		current = nested
	}
	return data
}

func TestObject(t *testing.T) {
	jsonTest := `
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

	rawJSONBytes := []byte(jsonTest)

	var jsonSerializedObject map[string]any
	var jsonSerializedObjectCopy map[string]any

	timeStartJSONDes := time.Now()
	err := json.Unmarshal(rawJSONBytes, &jsonSerializedObject)
	jsonDesDuration := time.Since(timeStartJSONDes)
	if err != nil {
		t.Fail()
	}

	err = json.Unmarshal(rawJSONBytes, &jsonSerializedObjectCopy)
	if err != nil {
		t.Fail()
	}

	timeStartUBON := time.Now()
	ubonBytes, err := ubonEncoder.Encode(jsonSerializedObject)
	durationSerializationUBON := time.Since(timeStartUBON)
	if err != nil {
		t.Fail()
	}

	timeStartUBONDes := time.Now()
	ubonRestoredRaw, err := ubonDecoder.Decode(ubonBytes)
	durationDeserializationUBON := time.Since(timeStartUBONDes)
	if err != nil {
		t.Fail()
	}
	ubonRestored := ubonRestoredRaw.(map[string]interface{})

	if !mapsEqual(ubonRestored, jsonSerializedObjectCopy) {
		t.Fail()
	}
	timeStartJSONSerail := time.Now()
	jsonRawData, err := json.Marshal(jsonSerializedObjectCopy)
	jsonSerializeDuration := time.Since(timeStartJSONSerail)
	if err != nil {
		t.Fail()
	}

	startCBORDeSerailization := time.Now()
	cborRawData, err := cbor.Marshal(jsonSerializedObjectCopy)
	desCBORDur := time.Since(startCBORDeSerailization)
	if err != nil {
		t.Fail()
	}

	var cborRestored map[string]interface{}
	startCBORSerailization := time.Now()
	err = cbor.Unmarshal(cborRawData, &cborRestored)
	cborSerializedDUr := time.Since(startCBORSerailization)
	if err != nil {
		log.Println("CBOR failed err: ", err)
	}

	ubonLen := float64(len(ubonBytes))
	jsonLen := float64(len(jsonRawData))
	cborLen := float64(len(cborRawData))
	profitJSON := float64(ubonLen / jsonLen)
	profitCBOR := float64(ubonLen / cborLen)

	//
	fmt.Println("cbor serialization duration: ", cborSerializedDUr)
	fmt.Println("json serialization duration: ", jsonSerializeDuration)
	fmt.Println("ubon serialization duration: ", durationSerializationUBON)

	fmt.Println("cbor deserialization duration: ", desCBORDur)
	fmt.Println("json deserialization duration: ", jsonDesDuration)
	fmt.Println("ubon deserialization duration: ", durationDeserializationUBON)

	fmt.Println("CBOR len: ", cborLen)
	fmt.Println("JSON len: ", jsonLen)
	fmt.Println("UBON len: ", ubonLen)
	fmt.Println("Profit relative to JSON: ", profitJSON)
	fmt.Println("Profit relative to CBOR: ", profitCBOR)
}

func mapsEqual(a, b map[string]interface{}) bool {
	return reflect.DeepEqual(a, b)
}
