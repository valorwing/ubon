package ubonTypes

const (
	// UBON type length 5 bit
	// base UBON Coding [ STRING_VAR_NAME_LEN | STRING_VAR_NAME | UBON_VAR_TYPE | DATA ]

	// Special
	UBON_nilOrNull = 0

	// Primitives (Sorted by size)
	UBON_Bool = 1 // 1 bit

	UBON_Int4  = 2 // 4 bit
	UBON_UInt4 = 3 // 4 bit

	UBON_Int8  = 4 // 8 bit
	UBON_UInt8 = 5 // 8 bit

	UBON_Int12  = 6 // 12 bit
	UBON_UInt12 = 7 // 12 bit

	UBON_Int16  = 8 // 16 bit
	UBON_UInt16 = 9 // 16 bit

	UBON_Int24  = 10 // 24 bit
	UBON_UInt24 = 11 // 24 bit

	UBON_Int32   = 12 // 32 bit
	UBON_UInt32  = 13 // 32 bit (also rune)
	UBON_Float32 = 14 // 32 bit

	UBON_Int56     = 15 // 56 bit
	UBON_UInt56    = 16 // 56 bit
	UBON_Float56   = 17 // 56 bit
	UBON_Complex56 = 18 // 56 bit

	UBON_Int64     = 19 // 64 bit
	UBON_UInt64    = 20 // 64 bit
	UBON_Float64   = 21 // 64 bit
	UBON_Complex64 = 22 // 64 bit

	UBON_Complex128 = 23 // 128 bit

	// Variative length (Strings, Arrays, Objects)
	UBON_RawData = 24 // raw data array
	UBON_String  = 25 // AKA char array

	UBON_Array  = 26
	UBON_Object = 27

	// Variative sized primitives
	UBON_VariativeSizedInt   = 28
	UBON_VariativeSizedUInt  = 29
	UBON_VariativeSizedFloat = 30

	// Special mark for using variative-sized primitives in arrays
	UBON_VariativeSizedArray = 31
)
