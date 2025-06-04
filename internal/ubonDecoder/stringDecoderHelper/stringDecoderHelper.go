package stringDecoderHelper

import (
	bitcode "ubon/internal/bitCode"
	bitcodeHashMap "ubon/internal/bitCodeHashMap"
	"ubon/internal/huffman"
	"ubon/internal/readOnlyBitStream"
)

func ReadEncodedString(rs *readOnlyBitStream.ReadOnlyBitStream, minCodeLen int, reverseCodes *bitcodeHashMap.BitcodeHashMap[rune]) (any, error) {

	outString := make([]rune, 0, 30)
	var currentBitCode bitcode.BitCode = bitcode.NewZeroBitCodeWithLength(0)
	isNewChar := true
	for {
		readedLen := 1
		if isNewChar {
			readedLen = minCodeLen
		}
		b, err := rs.ReadBitCode(readedLen)
		if err != nil {
			return nil, err
		}
		currentBitCode.AppendBitCode(*b)
		str, ok := reverseCodes.Get(currentBitCode)
		if !ok {
			isNewChar = false
			continue
		}
		currentBitCode.Clear()
		isNewChar = true
		if str == huffman.EOS_Char {
			break
		} else {
			outString = append(outString, str)
		}
	}
	return string(outString), nil
}
