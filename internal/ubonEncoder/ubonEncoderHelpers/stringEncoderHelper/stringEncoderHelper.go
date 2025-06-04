package stringEncoderHelper

import (
	"errors"
	bitcode "ubon/internal/bitCode"
	"ubon/internal/huffman"
	"ubon/internal/ubonOpCodes"
	"ubon/internal/writeOnlyBitStream"
)

func WriteEncodedStringToBitStream(input string, codes map[rune]bitcode.BitCode, ws *writeOnlyBitStream.WriteOnlyBitStream) error {
	ws.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_STRING)
	for _, v := range input {
		code, ok := codes[v]
		if !ok {
			return errors.New("writeEncodedStringToBitStream huffman broken")
		}
		ws.AppendBitCode(code)
	}
	ws.AppendBitCode(codes[huffman.EOS_Char])
	return nil
}

func WriteVarNameEncodedStringToBitStream(input string, codes map[rune]bitcode.BitCode, ws *writeOnlyBitStream.WriteOnlyBitStream) error {
	for _, v := range input {
		code, ok := codes[v]
		if !ok {
			return errors.New("writeEncodedStringToBitStream huffman broken")
		}
		ws.AppendBitCode(code)
	}
	ws.AppendBitCode(codes[huffman.EOS_Char])
	return nil
}
