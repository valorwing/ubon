package ubonTypes

const (

	//UBON type length 5 bit

	//Special

	// 1 bit
	UBON_nilOrNull = 0
	//Primitives

	// 1 bit
	UBON_Bool = 1

	// 8 bit
	UBON_Int8  = 2
	UBON_UInt8 = 3

	// 16 bit
	UBON_Int16  = 4
	UBON_UInt16 = 5

	// 32 bit
	UBON_Int32   = 6
	UBON_UInt32  = 7 // + rune
	UBON_Float32 = 8

	// 64 bit
	UBON_Int64     = 9
	UBON_UInt64    = 10
	UBON_Float64   = 11
	UBON_Complex64 = 12

	// 128 bit
	UBON_Complex128 = 13

	// Variative length
	UBON_Array  = 14
	UBON_Object = 15

	UBON_String = 16 // AKA char array

	// Unique half sized types self compressed type only internal used
	// for example input number int32 value 7 32 bit transformed
	// to UBON_Int4 value 7 4 bit and restored to int8 if no has layout
	// (target dictionary string -> any) or filled structre layout
	// in this example int32

	//4 bit
	UBON_Int4  = 17
	UBON_UInt4 = 18

	//12 bit
	UBON_Int12  = 19
	UBON_UInt12 = 20

	//16 bit half
	UBON_Float16 = 21

	//24 bit
	UBON_Int24  = 22
	UBON_UInt24 = 23

	// 56 bit
	UBON_Int56     = 24
	UBON_UInt56    = 25
	UBON_Float56   = 26
	UBON_Complex56 = 27

	// Variative sized section for example array with 50 elements int64 has max bit length 33
	// this sation can be writed as UBON_VariativeSizedArray | UBON_VariativeSizedInt | 33 |
	// | elements count | n1 , n2 ... |

	UBON_VariativeSizedInt   = 28
	UBON_VariativeSizedUInt  = 29
	UBON_VariativeSizedFloat = 30

	//Special mark for use varitive sized primitives
	UBON_VariativeSizedArray = 31
)
