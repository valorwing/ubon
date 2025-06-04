package booleanEncoderHelper

import (
	"ubon/internal/ubonOpCodes"
	"ubon/internal/writeOnlyBitStream"
)

func WriteBooleanToBitStream(value bool, ws *writeOnlyBitStream.WriteOnlyBitStream) {
	//writeOpCode and value
	ws.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_BOOLEAN.CloneAppend(value))
}
