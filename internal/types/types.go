package types

type UBON_Type uint8 // 5 bits reserved (max 32 types range 0-31 included)

// UBON Value types
const (
	UBON_Type_NULL    UBON_Type = 0 // Default zero value if primitive, or null/nil if complex
	UBON_Type_Boolean UBON_Type = 1

	// Encoded string (e.g. Huffman-compressed)
	UBON_Type_String UBON_Type = 2

	// Signed integers
	UBON_Type_Int8   UBON_Type = 3
	UBON_Type_Int16  UBON_Type = 4
	UBON_Type_Int32  UBON_Type = 5
	UBON_Type_Int64  UBON_Type = 6
	UBON_Type_Int128 UBON_Type = 7 // 128-bit signed integer

	// Unsigned integers
	UBON_Type_UInt8   UBON_Type = 8
	UBON_Type_UInt16  UBON_Type = 9
	UBON_Type_UInt32  UBON_Type = 10
	UBON_Type_UInt64  UBON_Type = 11
	UBON_Type_UInt128 UBON_Type = 12 // 128-bit unsigned integer

	// Raw bytes and characters
	UBON_Type_Byte        UBON_Type = 13 // 8-bit raw byte
	UBON_Type_Char_ASCII8 UBON_Type = 14 // 8-bit ASCII char (values 0-255) 7-bit table (ASCII7) included
	UBON_Type_Char_UTF8   UBON_Type = 15 // UTF-8 encoded character

	// Floating-point types (IEEE 754 or extended)
	UBON_Type_Float32   UBON_Type = 16 // 32-bit IEEE float
	UBON_Type_Float64   UBON_Type = 17 // 64-bit IEEE float aka double64
	UBON_Type_Double128 UBON_Type = 18 // 128-bit float (quad precision, IEEE 754-2008)

	//Others
	UBON_Type_UUID128     UBON_Type = 19 // 128-bit UUID
	UBON_Type_Timestamp64 UBON_Type = 20 // 64-bit timestamp
	UBON_Type_Decimal64   UBON_Type = 21 // 64-bit decimal
	UBON_Type_Decimal128  UBON_Type = 22 // 128-bit decimal
	UBON_Type_Bitfield    UBON_Type = 23 // variative length type after bitfield declaration - bitlen min 1 bit - max 256 bit

	UBON_Type_Object UBON_Type = 24 // aka key-value dictionary
	UBON_Type_Array  UBON_Type = 25

	//Reserve maybe support external pack in future
	UBON_Type_Custom_Pack1 UBON_Type = 26
	UBON_Type_Custom_Pack2 UBON_Type = 27
	UBON_Type_Custom_Pack3 UBON_Type = 28
	UBON_Type_Custom_Pack4 UBON_Type = 29
	UBON_Type_Custom_Pack5 UBON_Type = 30
	UBON_Type_Custom_Pack6 UBON_Type = 31
)
