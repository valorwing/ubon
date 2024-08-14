package ubonTypes

import (
	"reflect"
	"sync"
	"ubon/internal/errors"
	"unsafe"
)

const (
	CheckReject8To4BitsMask   = 0b00001111
	CheckReject16To12BitsMask = 0b0000111111111111
	CheckReject16To8BitsMask  = 0b0000000011111111
	CheckReject16To4BitsMask  = 0b0000000000001111
	CheckReject32To24BitsMask = 0b00000000111111111111111111111111
	CheckReject32To16BitsMask = 0b00000000000000001111111111111111
	CheckReject32To12BitsMask = 0b00000000000000000000111111111111
	CheckReject32To8BitsMask  = 0b00000000000000000000000011111111
	CheckReject32To4BitsMask  = 0b00000000000000000000000000001111
	//64 - 56
	//64 - 32
	//64 - 24
	//64 - 16
	//64 - 8
	//64 - 4
)

const (
	// UBON type length 5 bit
	// base UBON Coding [ STRING_VAR_NAME_LEN | STRING_VAR_NAME | UBON_VAR_TYPE | DATA ]

	// Special
	UBON_nilOrNull = uint8(0)

	// Primitives (Sorted by size)
	UBON_Bool = uint8(1) // 1 bit

	UBON_Int4  = uint8(2) // 4 bit
	UBON_UInt4 = uint8(3) // 4 bit

	UBON_Int8  = uint8(4) // 8 bit
	UBON_UInt8 = uint8(5) // 8 bit

	UBON_Int12  = uint8(6) // 12 bit
	UBON_UInt12 = uint8(7) // 12 bit

	UBON_Int16  = uint8(8) // 16 bit
	UBON_UInt16 = uint8(9) // 16 bit

	UBON_Int24  = uint8(10) // 24 bit
	UBON_UInt24 = uint8(11) // 24 bit

	UBON_Float64_28 = uint8(12) //  bit

	UBON_Int32   = uint8(13) // 32 bit
	UBON_UInt32  = uint8(14) // 32 bit (also rune)
	UBON_Float32 = uint8(15) // 32 bit

	UBON_Int56  = uint8(16) // 56 bit
	UBON_UInt56 = uint8(17) // 56 bit

	UBON_Int64     = uint8(18) // 64 bit
	UBON_UInt64    = uint8(19) // 64 bit
	UBON_Float64   = uint8(20) // 64 bit
	UBON_Complex64 = uint8(21) // 64 bit

	UBON_Complex128 = uint8(22) // 128 bit

	// Variative length (Strings, Arrays, Objects)
	UBON_String = uint8(23) // AKA char array

	UBON_RawData = uint8(24) // raw data array

	UBON_Array               = uint8(25)
	UBON_VariativeSizedArray = uint8(26)
	UBON_Object              = uint8(28)
)

var (
	exampleDefaultInt        = int(0)
	exampleDefaultUInt       = uint(0)
	defaultIntIs32BitsLen    = false
	defaultUIntIs32BitsLen   = false
	defaultIntLenDetectOnce  = sync.Once{}
	defaultUIntLenDetectOnce = sync.Once{}
)

func GetDefaultIntIs32BitsLen() bool {
	defaultIntLenDetectOnce.Do(func() {
		defaultIntIs32BitsLen = unsafe.Sizeof(exampleDefaultInt) == 4
	})
	return defaultIntIs32BitsLen
}

func GetDefaultUIntIs32BitsLen() bool {
	defaultUIntLenDetectOnce.Do(func() {
		defaultUIntIs32BitsLen = unsafe.Sizeof(exampleDefaultUInt) == 4
	})
	return defaultUIntIs32BitsLen
}

func DetectType(object any) (uint8, error) {

	reflectType := reflect.TypeOf(object)
	retVal := UBON_nilOrNull
	var err error = nil

	switch reflectType.Kind() {
	case reflect.Bool:
		retVal = UBON_Bool
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:

	case reflect.Uint:

	case reflect.Uint8:
		casted := object.(uint8)
		retVal = UBON_UInt8
		if (CheckReject8To4BitsMask & casted) == casted {
			retVal = UBON_UInt4
		}
	case reflect.Uint16:
		casted := object.(uint16)
		retVal = UBON_UInt16
		if (CheckReject16To4BitsMask & casted) == casted {
			retVal = UBON_UInt4
		} else if (CheckReject16To8BitsMask & casted) == casted {
			retVal = UBON_UInt8
		} else if (CheckReject16To12BitsMask & casted) == casted {
			retVal = UBON_UInt12
		}
	case reflect.Uint32:
		retVal = UBON_UInt32
		casted := object.(uint32)
		if (CheckReject32To4BitsMask & casted) == casted {
			retVal = UBON_UInt4
		} else if (CheckReject32To8BitsMask & casted) == casted {
			retVal = UBON_UInt8
		} else if (CheckReject32To12BitsMask & casted) == casted {
			retVal = UBON_UInt12
		} else if (CheckReject32To16BitsMask & casted) == casted {
			retVal = UBON_UInt16
		} else if (CheckReject32To24BitsMask & casted) == casted {
			retVal = UBON_UInt24
		}
	case reflect.Uint64:
		retVal = UBON_UInt64

	case reflect.Float32:
		retVal = UBON_Float32
	case reflect.Float64:
		retVal = UBON_Float64
	case reflect.Complex64:
		retVal = UBON_Complex64
	case reflect.Complex128:
		retVal = UBON_Complex128
	case reflect.Array:
	case reflect.Slice:
		retVal = UBON_Array
	case reflect.String:
		retVal = UBON_String
	case reflect.Map:
	case reflect.Struct:
		retVal = UBON_Object
	case reflect.Pointer:
		reflectValue := reflect.ValueOf(object)
		if reflectValue.IsNil() {
			return UBON_nilOrNull, nil
		}
		return DetectType(reflectValue.Elem().Interface())
	case reflect.UnsafePointer:
	case reflect.Uintptr:
	case reflect.Invalid:
	case reflect.Chan:
	case reflect.Func:
	case reflect.Interface:
		err = errors.UnsupportedCodingTypeError(reflectType.String() + " can't be serialized")
	}

	return retVal, err
}
