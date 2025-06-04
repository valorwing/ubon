package ubonOpCodes

import bitcode "ubon/internal/bitCode"

// TOP OPCODES 3 bits
const UBON_OP_CODE_BITLEN = 3

var (
	//000 - object coming
	UBON_OP_NEXT_OBJECT = bitcode.NewBitCodeWithBoolCode(false, false, false)
	//001 - array coming
	UBON_OP_NEXT_ARRAY = bitcode.NewBitCodeWithBoolCode(false, false, true)
	//010 - string coming
	UBON_OP_NEXT_STRING = bitcode.NewBitCodeWithBoolCode(false, true, false)
	//011 - int coming (sized after opcode sized mask 0 - 8, 10 - 16, 110 - 32, 111 - 64)
	UBON_OP_NEXT_INT = bitcode.NewBitCodeWithBoolCode(false, true, true)
	//100 - float coming (sized after opcode sized mask 0 - 32, 1 - 64)
	UBON_OP_NEXT_FLOAT = bitcode.NewBitCodeWithBoolCode(true, false, false)
	//101 - boolean coming
	UBON_OP_NEXT_BOOLEAN = bitcode.NewBitCodeWithBoolCode(true, false, true)
	//110 - uint coming
	UBON_OP_NEXT_UNSIGNED_INT = bitcode.NewBitCodeWithBoolCode(true, true, false)
	//111 - special opcode variative usage
	//name value space or single object opcode - null value
	//object space while await next object value type end of value
	//object space double special - end of file
	UBON_OP_NEXT_SPECIAL = bitcode.NewBitCodeWithBoolCode(true, true, true)
)
