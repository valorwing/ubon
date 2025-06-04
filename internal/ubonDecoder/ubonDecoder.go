package ubonDecoder

import (
	"errors"
	bitcode "ubon/internal/bitCode"
	bitcodeHashMap "ubon/internal/bitCodeHashMap"
	"ubon/internal/huffman"
	"ubon/internal/readOnlyBitStream"
	"ubon/internal/ubonDecoder/stringDecoderHelper"
	"ubon/internal/ubonDecoder/ubonDecoderHelpers/booleanDecoderHelper"
	"ubon/internal/ubonDecoder/ubonDecoderHelpers/floatDecoderHelper"
	"ubon/internal/ubonDecoder/ubonDecoderHelpers/intDecoderHelper"
	"ubon/internal/ubonDecoder/ubonDecoderHelpers/uIntDecoderHelper"
	"ubon/internal/ubonHeader"
	"ubon/internal/ubonOpCodes"
)

func Decode(input []byte) (any, error) {

	rs := readOnlyBitStream.NewReadOnlyBitStream(input)
	header, err := ubonHeader.ReadUbonHeaderFromReadOnlyBitStream(&rs)
	if err != nil {
		return nil, err
	}
	//if hasAlphabetReadAlphabet and decode tools
	var alphabet []rune
	var minHuffmanBitCodeLen int = 0
	var huffmanReverseCodesMap *bitcodeHashMap.BitcodeHashMap[rune]
	if header.AlphabetSectionIncluded {
		alphabet, err = huffman.AlphabetFromBitStream(&rs)
		if err != nil {
			return nil, err
		}
		tree := huffman.BuildTree(alphabet)
		if tree == nil {
			return nil, errors.New("huffman tree build failed")
		}
		codes := make(map[rune]bitcode.BitCode)
		huffman.GenerateCodes(tree, bitcode.NewZeroBitCodeWithLength(0), &codes)
		minBitCodeLen := -1
		reverseCodes := bitcodeHashMap.NewHashMap[rune](len(codes) / 2)
		for k, v := range codes {
			if minBitCodeLen == -1 {
				minBitCodeLen = v.BitLength()
			} else if minBitCodeLen > v.BitLength() {
				minBitCodeLen = v.BitLength()
			}
			reverseCodes.Put(v, k)
		}
		minHuffmanBitCodeLen = minBitCodeLen
		huffmanReverseCodesMap = reverseCodes
	}
	firstItemType, err := rs.ReadBitCode(3)
	if err != nil {
		return nil, err
	}
	//single primitive
	if firstItemType.Equal(ubonOpCodes.UBON_OP_NEXT_SPECIAL) {
		return nil, nil
	} else if firstItemType.Equal(ubonOpCodes.UBON_OP_NEXT_BOOLEAN) {
		return booleanDecoderHelper.ReadBool(&rs)
	} else if firstItemType.Equal(ubonOpCodes.UBON_OP_NEXT_INT) {
		return intDecoderHelper.ReadMaskedAutoInt(&rs)
	} else if firstItemType.Equal(ubonOpCodes.UBON_OP_NEXT_UNSIGNED_INT) {
		return uIntDecoderHelper.ReadMaskedAutoUnsignedInt(&rs)
	} else if firstItemType.Equal(ubonOpCodes.UBON_OP_NEXT_FLOAT) {
		return floatDecoderHelper.ReadMaskedAutoFloat(&rs)
	} else if firstItemType.Equal(ubonOpCodes.UBON_OP_NEXT_STRING) {
		if !header.AlphabetSectionIncluded {
			return nil, errors.New("single string alphabet required may be is not ubon frame?")
		} else {
			return stringDecoderHelper.ReadEncodedString(&rs, minHuffmanBitCodeLen, huffmanReverseCodesMap)
		}
	} else /*arrays and objects*/ if firstItemType.Equal(ubonOpCodes.UBON_OP_NEXT_OBJECT) {
		return decodeObject(&rs, huffmanReverseCodesMap, minHuffmanBitCodeLen)
	} else if firstItemType.Equal(ubonOpCodes.UBON_OP_NEXT_ARRAY) {
		panic("not implemented")
	}

	return nil, errors.New("failed read header or processing top opcode is ubon frame?")
}

func decodeNestedObject(rs *readOnlyBitStream.ReadOnlyBitStream, reverseCodes *bitcodeHashMap.BitcodeHashMap[rune], minCodeLen int) (any, error) {
	var nestedObject = map[string]interface{}{}
	for {
		opCode, err := rs.ReadBitCode(3)
		if err != nil {
			return nil, err
		}
		var currentVal any = nil

		//may be is nil or end brace
		if opCode.Equal(ubonOpCodes.UBON_OP_NEXT_SPECIAL) {
			flag, err := rs.ReadBitCode(1)
			if err != nil {
				return nil, err
			}
			if !flag.GetBit(0) {
				//zero nil value
				currentVal = nil
			} else {
				//enc object
				break
			}
		} else if opCode.Equal(ubonOpCodes.UBON_OP_NEXT_BOOLEAN) {
			currentVal, err = booleanDecoderHelper.ReadBool(rs)
			if err != nil {
				return nil, err
			}
		} else if opCode.Equal(ubonOpCodes.UBON_OP_NEXT_INT) {
			intValRaw, err := intDecoderHelper.ReadMaskedAutoInt(rs)
			if err != nil {
				return nil, err
			}
			switch intValTyped := intValRaw.(type) {
			case int8:
				currentVal = float64(intValTyped)
			case int16:
				currentVal = float64(intValTyped)
			case int:
				currentVal = float64(intValTyped)
			case int32:
				currentVal = float64(intValTyped)
			case int64:
				currentVal = float64(intValTyped)
			default:
				return nil, errors.New("broken int type")
			}
		} else if opCode.Equal(ubonOpCodes.UBON_OP_NEXT_UNSIGNED_INT) {
			uIntValRaw, err := uIntDecoderHelper.ReadMaskedAutoUnsignedInt(rs)
			if err != nil {
				return nil, err
			}
			switch uIntValTyped := uIntValRaw.(type) {
			case uint8:
				currentVal = float64(uIntValTyped)
			case uint16:
				currentVal = float64(uIntValTyped)
			case uint:
				currentVal = float64(uIntValTyped)
			case uint32:
				currentVal = float64(uIntValTyped)
			case uint64:
				currentVal = float64(uIntValTyped)
			default:
				return nil, errors.New("broken uint type")
			}
		} else if opCode.Equal(ubonOpCodes.UBON_OP_NEXT_FLOAT) {
			floatRawVal, err := floatDecoderHelper.ReadMaskedAutoFloat(rs)
			if err != nil {
				return nil, err
			}
			switch floatTypedVal := floatRawVal.(type) {
			case float32:
				currentVal = float64(floatTypedVal)
			case float64:
				currentVal = floatTypedVal
			}
		} else if opCode.Equal(ubonOpCodes.UBON_OP_NEXT_STRING) {
			currentVal, err = stringDecoderHelper.ReadEncodedString(rs, minCodeLen, reverseCodes)
			if err != nil {
				return nil, err
			}
		} else if opCode.Equal(ubonOpCodes.UBON_OP_NEXT_OBJECT) {
			currentVal, err = decodeNestedObject(rs, reverseCodes, minCodeLen)
			if err != nil {
				return nil, err
			}
		}

		currentKey, err := stringDecoderHelper.ReadEncodedString(rs, minCodeLen, reverseCodes)
		if err != nil {
			return nil, err
		}
		currentKeyString := currentKey.(string)
		nestedObject[currentKeyString] = currentVal
	}
	return nestedObject, nil
}

func decodeObject(rs *readOnlyBitStream.ReadOnlyBitStream, reverseCodes *bitcodeHashMap.BitcodeHashMap[rune], minCodeLen int) (any, error) {

	var rootObject = map[string]interface{}{}
	for {
		opCode, err := rs.ReadBitCode(3)
		if err != nil {
			return nil, err
		}
		var currentVal any = nil

		//may be is nil or end brace
		if opCode.Equal(ubonOpCodes.UBON_OP_NEXT_SPECIAL) {
			flag, err := rs.ReadBitCode(1)
			if err != nil {
				return nil, err
			}
			if !flag.GetBit(0) {
				//zero nil value
				currentVal = nil
			} else {
				//enc object
				break
			}
		} else if opCode.Equal(ubonOpCodes.UBON_OP_NEXT_BOOLEAN) {
			currentVal, err = booleanDecoderHelper.ReadBool(rs)
			if err != nil {
				return nil, err
			}
		} else if opCode.Equal(ubonOpCodes.UBON_OP_NEXT_INT) {
			intValRaw, err := intDecoderHelper.ReadMaskedAutoInt(rs)
			if err != nil {
				return nil, err
			}
			switch intValTyped := intValRaw.(type) {
			case int8:
				currentVal = float64(intValTyped)
			case int16:
				currentVal = float64(intValTyped)
			case int:
				currentVal = float64(intValTyped)
			case int32:
				currentVal = float64(intValTyped)
			case int64:
				currentVal = float64(intValTyped)
			default:
				return nil, errors.New("broken int type")
			}
		} else if opCode.Equal(ubonOpCodes.UBON_OP_NEXT_UNSIGNED_INT) {
			uIntValRaw, err := uIntDecoderHelper.ReadMaskedAutoUnsignedInt(rs)
			if err != nil {
				return nil, err
			}
			switch uIntValTyped := uIntValRaw.(type) {
			case uint8:
				currentVal = float64(uIntValTyped)
			case uint16:
				currentVal = float64(uIntValTyped)
			case uint:
				currentVal = float64(uIntValTyped)
			case uint32:
				currentVal = float64(uIntValTyped)
			case uint64:
				currentVal = float64(uIntValTyped)
			default:
				return nil, errors.New("broken uint type")
			}
		} else if opCode.Equal(ubonOpCodes.UBON_OP_NEXT_FLOAT) {
			floatRawVal, err := floatDecoderHelper.ReadMaskedAutoFloat(rs)
			if err != nil {
				return nil, err
			}
			switch floatTypedVal := floatRawVal.(type) {
			case float32:
				currentVal = float64(floatTypedVal)
			case float64:
				currentVal = floatTypedVal
			}
		} else if opCode.Equal(ubonOpCodes.UBON_OP_NEXT_STRING) {
			currentVal, err = stringDecoderHelper.ReadEncodedString(rs, minCodeLen, reverseCodes)
			if err != nil {
				return nil, err
			}
		} else if opCode.Equal(ubonOpCodes.UBON_OP_NEXT_OBJECT) {
			currentVal, err = decodeNestedObject(rs, reverseCodes, minCodeLen)
			if err != nil {
				return nil, err
			}
		}

		currentKey, err := stringDecoderHelper.ReadEncodedString(rs, minCodeLen, reverseCodes)
		if err != nil {
			return nil, err
		}
		currentKeyString := currentKey.(string)
		rootObject[currentKeyString] = currentVal
	}
	return rootObject, nil
}
