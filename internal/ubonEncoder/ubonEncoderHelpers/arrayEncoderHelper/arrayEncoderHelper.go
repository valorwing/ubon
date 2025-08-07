package arrayEncoderHelper

import (
	"reflect"
	bitcode "ubon/internal/bitCode"
	"ubon/internal/huffman"
	"ubon/internal/ubonOpCodes"
)

type ArrayCheckResult struct {
	IsUBONCompatible          bool
	IsAlphabetRequired        bool
	IsPrimitiveArray          bool
	IsOptionalFeatureRequried bool
	IsEmptyArray              bool
	ArrayContract             bitcode.BitCode
}

// UBON_OP_NEXT_ARRAY NEXT_ARRAY|[ARRAY_LEN_PREFIX][ARRAY_LEN][ARRAY_CONTRACT]
func AnalyzeArray(array []any, freqMap *huffman.HuffmanStringFrequencyMap) ArrayCheckResult {

	var retVal ArrayCheckResult = ArrayCheckResult{}
	if array == nil || len(array) == 0 {
		retVal.IsEmptyArray = true
		retVal.IsUBONCompatible = true
		retVal.IsAlphabetRequired = false
		retVal.ArrayContract = bitcode.NewZeroBitCodeWithLength(0)
		//UBON_OP_NEXT_ARRAY NEXT_ARRAY
		retVal.ArrayContract.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_ARRAY)
		//ARRAY_LEN_PREFIX 0 (0 - 8 bits, 10 - 16 bits, 110 - 32 bits, 111 - 64 bits)
		retVal.ArrayContract.Append(false)
		//ARRAY_LEN 8 bits zero (00000000)
		retVal.ArrayContract.AppendBitCode(bitcode.NewBitCodeFromBytes(0))

		return retVal
	} else if t, ok := checkArrayElementsType(array); ok {
		//multy dimensional array
		if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
			retVal.IsUBONCompatible = true
		} else {

		}

	} else {
		retVal.IsUBONCompatible = false
		return retVal
	}
}

// retval array type and bool flag is array is UBON compatible
func checkArrayElementsType(arr []any) (reflect.Type, bool) {
	if len(arr) == 0 {
		return nil, true
	}
	t := reflect.TypeOf(arr[0])

	if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
		for i := 1; i < len(arr); i++ {
			arrayElemType := reflect.TypeOf(arr[i])
			if arrayElemType.Kind() != reflect.Array && arrayElemType.Kind() != reflect.Slice {
				return nil, false
			}
		}
	} else {

		for i := 1; i < len(arr); i++ {
			if reflect.TypeOf(arr[i]) != t {
				return nil, false
			}
		}
		return t, true
	}
	return nil, false
}
