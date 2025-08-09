package rootElementHelper

import (
	"reflect"
	bitcode "ubon/internal/bitCode"
	"ubon/internal/ubonOpCodes"
)


//NEXT_OBJECT [OBJECT_CONTRACT][END_OBJECT_CONTRACT][OBJECT_DATA] HUFFNAN_OBJECT_NAME HUFFMAN EOS (root object unnamed)
//NEXT_ARRAY [ARRAY_PREFIXED_LENGTH][ARRAY_CONTRACT][ARRAY_DATA] HUFFMAN_ARRAY_NAME HUFFMAN EOS (root array unnamed)
//NEXT_PRIMITIVE_TYPE(AKA NEXT_BOOL NEXT_INT etc) [PRIMITIVE_DATA] HUFFMAN_PRIMITIVE_NAME HUFFMAN EOS (root primitive unnamed)
//NEXT_SPECIAL flag only if root nil

type RootElementCheckResult struct {
	IsNull                 bool
	IsPrimitive            bool
	IsObject               bool
	IsEmptyArray           bool
	IsPrimitiveArray       bool
	IsObjectArray          bool
	IsMutidimensionalArray bool
	Contract               bitcode.BitCode
}

func AnalyzeRootElementAndConstructPartialContract(rootElement any) (RootElementCheckResult, error) {
	var retVal RootElementCheckResult
	retVal.Contract = bitcode.NewZeroBitCodeWithLength(0)

	if rootElement == nil {
		retVal.IsNull = true
		retVal.Contract.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_SPECIAL)
		return retVal, nil
	}

	switch reflect.TypeOf(rootElement).Kind() {
	case reflect.Bool:
		retVal.IsPrimitive = true
		retVal.Contract.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_BOOL)

	reflectType := reflect.TypeOf(rootElement)
	reflect.Kind = reflectType.Kind()
	
	switch reflect.Kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		retVal.IsPrimitive = true
		retVal.Contract.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_INT)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		retVal.IsPrimitive = true
		retVal.Contract.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_UNSIGNED_INT)
	case reflect.Float32, reflect.Float64:
		retVal.IsPrimitive = true
		retVal.Contract.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_FLOAT)
	case reflect.String:
		retVal.IsPrimitive = true
		retVal.Contract.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_STRING)
	case reflect.Array,
}
}
