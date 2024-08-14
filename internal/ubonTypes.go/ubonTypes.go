package ubonTypes

import (
	"reflect"
	"sync"
	"ubon/internal/errors"
	"unsafe"
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
	CheckReject64To56BitsMask = 0b0000000011111111111111111111111111111111111111111111111111111111
	CheckReject64To32BitsMask = 0b0000000000000000000000000000000011111111111111111111111111111111
	CheckReject64To24BitsMask = 0b0000000000000000000000000000000000000000111111111111111111111111
	CheckReject64To16BitsMask = 0b0000000000000000000000000000000000000000000000001111111111111111
	CheckReject64To8BitsMask  = 0b0000000000000000000000000000000000000000000000000000000011111111
	CheckReject64To4BitsMask  = 0b0000000000000000000000000000000000000000000000000000000000001111
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
	kind := reflectType.Kind()

	switch {
	case kind == reflect.Bool:
		retVal = UBON_Bool
	case kind == reflect.Int8:
		casted := object.(int8)
		retVal = UBON_Int8
		if (CheckReject8To4BitsMask & casted) == casted {
			retVal = UBON_Int4
		}
	case kind == reflect.Int16:
		casted := object.(int16)
		retVal = UBON_Int16
		if (CheckReject16To4BitsMask & casted) == casted {
			retVal = UBON_Int4
		} else if (CheckReject16To8BitsMask & casted) == casted {
			retVal = UBON_Int8
		} else if (CheckReject16To12BitsMask & casted) == casted {
			retVal = UBON_Int12
		}
	case kind == reflect.Int32 || (kind == reflect.Int && GetDefaultIntIs32BitsLen()):
		retVal = UBON_UInt32
		var casted int32
		if kind == reflect.Uint {
			casted = int32(object.(uint))
		} else {
			casted = object.(int32)
		}
		if (CheckReject32To4BitsMask & casted) == casted {
			retVal = UBON_Int4
		} else if (CheckReject32To8BitsMask & casted) == casted {
			retVal = UBON_Int8
		} else if (CheckReject32To12BitsMask & casted) == casted {
			retVal = UBON_Int12
		} else if (CheckReject32To16BitsMask & casted) == casted {
			retVal = UBON_Int16
		} else if (CheckReject32To24BitsMask & casted) == casted {
			retVal = UBON_Int24
		}
	case kind == reflect.Int64 || (kind == reflect.Int && !GetDefaultIntIs32BitsLen()):
		retVal = UBON_Int64
		var casted int64
		if kind == reflect.Int {
			casted = int64(object.(int))
		} else {
			casted = object.(int64)
		}

		if (CheckReject64To4BitsMask & casted) == casted {
			retVal = UBON_Int4
		} else if (CheckReject64To8BitsMask & casted) == casted {
			retVal = UBON_Int8
		} else if (CheckReject64To16BitsMask & casted) == casted {
			retVal = UBON_Int16
		} else if (CheckReject64To24BitsMask & casted) == casted {
			retVal = UBON_Int24
		} else if (CheckReject64To32BitsMask & casted) == casted {
			retVal = UBON_Int32
		} else if (CheckReject64To56BitsMask & casted) == casted {
			retVal = UBON_Int56
		}

	case kind == reflect.Uint8:
		casted := object.(uint8)
		retVal = UBON_UInt8
		if (CheckReject8To4BitsMask & casted) == casted {
			retVal = UBON_UInt4
		}
	case kind == reflect.Uint16:
		casted := object.(uint16)
		retVal = UBON_UInt16
		if (CheckReject16To4BitsMask & casted) == casted {
			retVal = UBON_UInt4
		} else if (CheckReject16To8BitsMask & casted) == casted {
			retVal = UBON_UInt8
		} else if (CheckReject16To12BitsMask & casted) == casted {
			retVal = UBON_UInt12
		}
	case kind == reflect.Uint32 || (kind == reflect.Uint && GetDefaultUIntIs32BitsLen()):
		retVal = UBON_UInt32
		var casted uint32
		if kind == reflect.Uint {
			casted = uint32(object.(uint))
		} else {
			casted = object.(uint32)
		}

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
	case kind == reflect.Uint64 || (kind == reflect.Uint && !GetDefaultUIntIs32BitsLen()):
		retVal = UBON_UInt64
		var casted uint64
		if kind == reflect.Uint {
			casted = uint64(object.(uint))
		} else {
			casted = object.(uint64)
		}

		if (CheckReject64To4BitsMask & casted) == casted {
			retVal = UBON_UInt4
		} else if (CheckReject64To8BitsMask & casted) == casted {
			retVal = UBON_UInt8
		} else if (CheckReject64To16BitsMask & casted) == casted {
			retVal = UBON_UInt16
		} else if (CheckReject64To24BitsMask & casted) == casted {
			retVal = UBON_UInt24
		} else if (CheckReject64To32BitsMask & casted) == casted {
			retVal = UBON_UInt32
		} else if (CheckReject64To56BitsMask & casted) == casted {
			retVal = UBON_UInt56
		}

	case kind == reflect.Float32:
		retVal = UBON_Float32
	case kind == reflect.Float64:
		retVal = UBON_Float64
	case kind == reflect.Complex64:
		retVal = UBON_Complex64
	case kind == reflect.Complex128:
		retVal = UBON_Complex128
	case kind == reflect.Slice || kind == reflect.Array:
		retVal = UBON_Array
	case kind == reflect.String:
		retVal = UBON_String
	case kind == reflect.Struct || kind == reflect.Map:
		retVal = UBON_Object
	case kind == reflect.Pointer:
		reflectValue := reflect.ValueOf(object)
		if reflectValue.IsNil() {
			return UBON_nilOrNull, nil
		}
		return DetectType(reflectValue.Elem().Interface())
	default:
		err = errors.UnsupportedCodingTypeError(reflectType.String() + " can't be serialized")
	}

	return retVal, err
}
