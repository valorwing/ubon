package huffman

import (
	"errors"
	bitcode "ubon/internal/bitCode"
	"ubon/internal/readOnlyBitStream"
)

const EOS_Char string = "\u001A"

var EOS_Char_Bitcode bitcode.BitCode = bitcode.NewBitCodeWithBoolCode(false, false, false, true, true, false, true, false)

func AlphabetToBitcode(alphabet []string) (*bitcode.BitCode, error) {
	retValBitcode := bitcode.NewZeroBitCodeWithLength(0)
	eosCharCount := 0
	for _, s := range alphabet {
		if s == EOS_Char {
			eosCharCount++
		}
		data := []byte(s)
		retValBitcode.AppendBitCode(bitcode.NewBitCodeFromBytes(data...))
	}
	//single EOS_Char Required
	if eosCharCount != 1 {
		if eosCharCount == 0 {
			return nil, errors.New("in the alphabet according to the protocol it is obligatory to mention the special separator \u001A")
		} else {
			return nil, errors.New("alphabet serialization protocol error. The \u001A character may be mentioned in a string or variable name")
		}
	}
	retValBitcode.AppendBitCode(EOS_Char_Bitcode)
	return &retValBitcode, nil
}

// AlphabetFromBitStream reads a serialized alphabet from the bit stream.
// The EOS_Char (\u001A) plays a dual role:
//  1. It must be present exactly once in the alphabet as a valid symbol.
//  2. Its second occurrence marks the end of the alphabet and the beginning of payload.
//
// This ensures fixed-width 8-bit alignment for symbols without requiring a separate length prefix.
func AlphabetFromBitStream(bs *readOnlyBitStream.ReadOnlyBitStream) ([]string, error) {
	findFirstEOS := false
	alphabetBytes := make([]byte, 0)
	for {
		v, err := bs.ReadBitCode(8)
		if err != nil {
			return nil, err
		}
		if v.Equal(EOS_Char_Bitcode) {
			if findFirstEOS {
				break
			} else {
				findFirstEOS = true
				alphabetBytes = append(alphabetBytes, v.Bytes()...)
			}
		} else {
			alphabetBytes = append(alphabetBytes, v.Bytes()...)
		}

	}
	alphabetString := string(alphabetBytes)
	alphabet := make([]string, 0, len(alphabetString))
	for _, s := range alphabetString {
		alphabet = append(alphabet, string(s))
	}
	return alphabet, nil
}
