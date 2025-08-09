package arrayEncoderHelper

import (
	"reflect"
	bitcode "ubon/internal/bitCode"
	"ubon/internal/ubonOpCodes"
)

type ArrayCheckResult struct {
	IsUBONCompatible       bool
	IsPrimitiveArray       bool
	IsMutidimensionalArray bool
	IsEmptyArray           bool
	ArrayContract          bitcode.BitCode
}

func PrefixedArraySizeBitCode(arrayLen uint64) bitcode.BitCode {
	// ARRAY_LEN_PREFIX 0 (0 - 8 bits, 10 - 16 bits, 110 - 32 bits, 111 - 64 bits) and ARRAY_LEN
	retVal := bitcode.NewZeroBitCodeWithLength(0)
	if arrayLen <= 0b11111111 { // 8 bits
		retVal.Append(false) // 0
		retVal.AppendBitCode(bitcode.NewBitCodeFromBytes(byte(arrayLen)))
	} else if arrayLen <= 0b1111111111111111 { // 16 bits
		retVal.Append(true)  // 1
		retVal.Append(false) // 0
		retVal.AppendBitCode(bitcode.NewBitCodeFromBytes(byte(arrayLen), byte(arrayLen>>8)))
		//Write
	} else if arrayLen <= 0b11111111111111111111111111111111 { // 32 bits
		retVal.Append(true)  // 1
		retVal.Append(true)  // 1
		retVal.Append(false) // 0
		retVal.AppendBitCode(bitcode.NewBitCodeFromBytes(
			byte(arrayLen), byte(arrayLen>>8),
			byte(arrayLen>>16), byte(arrayLen>>24)))
	} else { // 64 bits
		retVal.Append(true) // 1
		retVal.Append(true) // 1
		retVal.Append(true) // 1
		retVal.AppendBitCode(bitcode.NewBitCodeFromBytes(
			byte(arrayLen), byte(arrayLen>>8),
			byte(arrayLen>>16), byte(arrayLen>>24),
			byte(arrayLen>>32), byte(arrayLen>>40),
			byte(arrayLen>>48), byte(arrayLen>>56)))
	}
	return retVal
}

// Partial contract for array with size and type of elements
func AnalyzeArrayAndContructPartialArrayContract(array []any) ArrayCheckResult {

	var retVal ArrayCheckResult = ArrayCheckResult{}
	retVal.ArrayContract = bitcode.NewZeroBitCodeWithLength(0)
	if len(array) == 0 || array == nil {
		retVal.IsEmptyArray = true
		retVal.IsUBONCompatible = true

		retVal.ArrayContract.AppendBitCode(PrefixedArraySizeBitCode(0))

		return retVal
	} else {

		arrayLen := uint64(len(array))
		retVal.ArrayContract.AppendBitCode(PrefixedArraySizeBitCode(arrayLen))
		var rType reflect.Type
		rType, retVal.IsUBONCompatible = checkArrayElementsType(array)
		if !retVal.IsUBONCompatible {
			return retVal
		}
		rKind := rType.Kind()
		switch rKind {
		case reflect.Array, reflect.Slice:
			//each nested array has its own contract
			retVal.IsMutidimensionalArray = true
			retVal.ArrayContract.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_ARRAY)
			return retVal
		case reflect.Bool:
			retVal.IsPrimitiveArray = true
			retVal.ArrayContract.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_BOOLEAN)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			retVal.IsPrimitiveArray = true
			retVal.IsUBONCompatible = true
			//aka auto int with individual bitsize prefix
			retVal.ArrayContract.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_INT)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			retVal.IsUBONCompatible = true
			retVal.IsPrimitiveArray = true
			//aka auto uint with individual bitsize prefix
			retVal.ArrayContract.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_UNSIGNED_INT)
		case reflect.Float32, reflect.Float64:
			retVal.IsUBONCompatible = true
			retVal.IsPrimitiveArray = true
			//aka auto float with individual bitsize prefix
			retVal.ArrayContract.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_FLOAT)
		case reflect.String:
			retVal.IsUBONCompatible = true
			retVal.IsPrimitiveArray = true
			//string with alphabet
			retVal.ArrayContract.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_STRING)
		case reflect.Map:
			retVal.IsUBONCompatible = true
			retVal.ArrayContract.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_OBJECT)

			return retVal
			//partial contract can detect this case later in external code
		case reflect.Struct:
			retVal.IsUBONCompatible = true
			retVal.ArrayContract.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_OBJECT)

			return retVal
			//partial contract can detect this case later in external code
		case reflect.Interface:
			retVal.IsUBONCompatible = true
			retVal.ArrayContract.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_OBJECT)

			return retVal
			//partial contract can detect this case later in external code
		case reflect.Pointer:
			//unwrap pointer and again check
			newArray := make([]any, len(array))
			for i, v := range array {
				if v != nil {
					newArray[i] = v.(*interface{})
				} else {
					newArray[i] = nil
				}
			}
			return AnalyzeArrayAndContructPartialArrayContract(newArray)
		default:
			retVal.IsUBONCompatible = false
			return retVal
		}
	}
	return retVal
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
