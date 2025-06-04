package ubonEncoder

import (
	"encoding/json"
	"errors"
	bitcode "ubon/internal/bitCode"
	"ubon/internal/huffman"
	"ubon/internal/ubonEncoder/ubonEncoderHelpers/booleanEncoderHelper"
	"ubon/internal/ubonEncoder/ubonEncoderHelpers/floatEncoderHelper"
	"ubon/internal/ubonEncoder/ubonEncoderHelpers/intEncoderHelper"
	"ubon/internal/ubonEncoder/ubonEncoderHelpers/stringEncoderHelper"
	uIntEncoderHelper "ubon/internal/ubonEncoder/ubonEncoderHelpers/uintEncoderHelper"
	"ubon/internal/ubonHeader"
	"ubon/internal/ubonOpCodes"
	"ubon/internal/writeOnlyBitStream"
)

// Prototype!!!! Allow any - json compitable object
// input allowed map[string]any, []any, or single null, bool, int, float, string value
func Encode(input any) ([]byte, error) {

	checkedInput, err := inputAlignToJSONCompitableIfNeeded(input)
	if err != nil {
		return nil, err
	}
	if input == nil {
		return encodeNullFirst()
	}
	switch checkedInput.(type) {
	case map[string]interface{}:
		return encodeObjectFirst(input.(map[string]interface{}))
	case bool:
		return encodeBooleanFirst(input.(bool))
	case []any:
		return encodeArrayFirst(input.([]interface{}))
	case int, int8, int16, int32, int64:
		return encodeIntegerFirst(input)
	case uint, uint8, uint16, uint32, uint64:
		return encodeUnsignedIntegerFirst(input)
	case float32, float64:
		return encodeFloatFirst(input)
	case string:
		return encodeStringFirst(input)
	}
	return nil, errors.New("undetected/unsupported type")
}

func inputAlignToJSONCompitableIfNeeded(input any) (any, error) {
	switch input.(type) {
	case
		bool,
		map[string]interface{},
		[]any,
		string,
		float32,
		float64,
		int8,
		int16,
		int32,
		int64,
		int,
		uint8,
		uint16,
		uint32,
		uint64,
		uint:
		return input, nil
	case uintptr, complex64, complex128:
		return nil, errors.New("unsupported type")
	default:
		//try temp support custom types
		var i map[string]interface{}
		b, err := json.Marshal(input)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, &i)
		if err != nil {
			return nil, err
		}
		return inputAlignToJSONCompitableIfNeeded(i)
	}
}

func recursiveCollectStringStats(input map[string]interface{}, freqMap *huffman.HuffmanStringFrequencyMap) error {

	for k, v := range input {
		freqMap.SendString(k)
		switch vTyped := v.(type) {
		case string:
			freqMap.SendString(vTyped)
		case map[string]interface{}:
			err := recursiveCollectStringStats(vTyped, freqMap)
			if err != nil {
				return err
			}
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64,
			float32, float64, bool, nil:
			continue
		case []any:
			panic("TODO://ARRAYS FEATURE")
		default:
			return errors.New("recursiveCollectStringStats find unsupported type")
		}
	}
	return nil
}

// map[string]any only //root object //TODO ARRAYS
func encodeObjectFirst(input map[string]interface{}) ([]byte, error) {
	if input == nil {
		//SPECIAL CASE
		return encodeNullFirst()
	}

	//varnames exist always has alphabed
	header := ubonHeader.NewDefaultUbonHeader(true)
	header.AlphabetSectionIncluded = true
	//first step huffman
	freqMap := huffman.NewHuffmanStringFrequencyMap()
	err := recursiveCollectStringStats(input, freqMap)
	if err != nil {
		return nil, err
	}
	freqMap.FinishСollectingStrings()

	ws := writeOnlyBitStream.NewWriteOnlyBitStream()
	ws.AppendBitCode(header.BitCode())
	alphabet := freqMap.GetAlphabet()
	alphabetBitcode, err := huffman.AlphabetToBitcode(alphabet)
	if err != nil {
		return nil, err
	}
	ws.AppendBitCode(*alphabetBitcode)
	tree := huffman.BuildTree(alphabet)
	if tree == nil {

		return nil, errors.New("huffman tree build failed")
	}
	codes := make(map[rune]bitcode.BitCode)
	huffman.GenerateCodes(tree, bitcode.NewZeroBitCodeWithLength(0), &codes)

	//ALPHABET READY OPEN ROOT OBJECT BRACE
	ws.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_OBJECT)
	//OBJECT CONTEXT SPECIAL OPCODE + 0 - null value
	//OBJECT CONTEXT SPECIAL OPCODE + 1 - object close brace

	//OPCODE:VALUE:VARNAME
	for k, vRaw := range input {
		switch v := vRaw.(type) {
		case nil:
			//OBJECT CONTEXT SPECIAL OPCODE + 0 - null value
			ws.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_SPECIAL)
			ws.AppendBitCode(bitcode.NewBitCodeWithBoolCode(false))

		case bool:
			booleanEncoderHelper.WriteBooleanToBitStream(v, ws)
		case int, int8, int16, int32, int64:
			err := intEncoderHelper.WriteAutoIntToBitStream(v, ws)
			if err != nil {
				return nil, err
			}
		case uint, uint8, uint16, uint32, uint64:
			err := uIntEncoderHelper.WriteAutoUnsignedIntToBitStream(v, ws)
			if err != nil {
				return nil, err
			}
		case float32:
			err := floatEncoderHelper.WriteAutoFloatToBitStream(v, ws)
			if err != nil {
				return nil, err
			}
		case float64:
			err := floatEncoderHelper.WriteAutoFloatToBitStream(v, ws)
			if err != nil {
				return nil, err
			}
		case string:
			err := stringEncoderHelper.WriteEncodedStringToBitStream(v, codes, ws)
			if err != nil {
				return nil, err
			}

		case map[string]interface{}:
			//nested object
			err := encodeNestedObject(v, ws, codes)
			if err != nil {
				return nil, err
			}
		case []any:
			//nested array
			panic("TODO://ARRAYS FEATURE")
		}

		err := stringEncoderHelper.WriteVarNameEncodedStringToBitStream(k, codes, ws)
		if err != nil {
			return nil, err
		}
	}
	//CLOSE ROOT BRACE ubonOpCodes.UBON_OP_NEXT_SPECIAL+1
	ws.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_SPECIAL.CloneAppend(true))
	return ws.Bytes(), nil
}

func encodeNestedObject(input map[string]interface{}, ws *writeOnlyBitStream.WriteOnlyBitStream, huffmanCodes map[rune]bitcode.BitCode) error {
	ws.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_OBJECT)

	for k, vRaw := range input {
		switch v := vRaw.(type) {
		case nil:
			//OBJECT CONTEXT SPECIAL OPCODE + 0 - null value
			ws.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_SPECIAL)
			ws.AppendBitCode(bitcode.NewBitCodeWithBoolCode(false))

		case bool:
			booleanEncoderHelper.WriteBooleanToBitStream(v, ws)
		case int, int8, int16, int32, int64:
			err := intEncoderHelper.WriteAutoIntToBitStream(v, ws)
			if err != nil {
				return err
			}
		case uint, uint8, uint16, uint32, uint64:
			err := uIntEncoderHelper.WriteAutoUnsignedIntToBitStream(v, ws)
			if err != nil {
				return err
			}
		case float32:
			err := floatEncoderHelper.WriteAutoFloatToBitStream(v, ws)
			if err != nil {
				return err
			}
		case float64:
			err := floatEncoderHelper.WriteAutoFloatToBitStream(v, ws)
			if err != nil {
				return err
			}

		case map[string]interface{}:
			//nested object
			err := encodeNestedObject(v, ws, huffmanCodes)
			if err != nil {
				return err
			}
		case string:
			err := stringEncoderHelper.WriteEncodedStringToBitStream(v, huffmanCodes, ws)
			if err != nil {
				return err
			}
		case []any:
			//nested array
			panic("TODO://ARRAYS FEATURE")
		}

		err := stringEncoderHelper.WriteVarNameEncodedStringToBitStream(k, huffmanCodes, ws)
		if err != nil {
			return err
		}
	}

	//CLOSE BRACE
	ws.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_SPECIAL.CloneAppend(true))
	return nil
}

// []any

//ARRAYS
//ARRAYS OF PRIMITIVE
//	[OP_CODE_ARRAY][OP_CODE_PRIMITIVE][ARRAY_LEN_PREFIX][ARRAY_LEN_UINT][PRIMITIVES...ARRAY_LEN]
//ARRAYS OF OBJECTS
// [OP_CODE_ARRAY][OP_CODE_OBJECT][OP_CODE|NAME object_fields_pair 3bits + EOS string while reached OP_CODE_SPECIAL]
// [OPTIONAL_MASK_FEATURE flag 1 bit][all optinal functional not writed if OPTIONAL_FLAG == 0 OPTIONAL_MASK - bits array len object_fields_pair and indicated positive flags for option mask object example][ARRAY_LEN_PREFIX][ARRAY_LEN_UINT] | [OBJECT EXAMPLE OPTION MASK (only enabled) and OBJECT_DATA ordered by object_fields_pair] * ARRAY_LEN|
// MULTIDIMENTIONAL
// nested frame description like [OP_CODE_ARRAY] open brace [OP_CODE_SPECIAL] close brace
//

func encodeArrayFirst(input []interface{}) ([]byte, error) {

	ws := writeOnlyBitStream.NewWriteOnlyBitStream()
	isAlphabetRequired := false
	//TODO: implement alphabet detection primitive arrays no need alphabet
	header := ubonHeader.NewDefaultUbonHeader(true)
	//TODO: implement array
	//stub
	_, _, _ = header, isAlphabetRequired, ws
	return nil, nil
}

// int, string, float, boolean, null !!!Unnamed
func encodeIntegerFirst(input any) ([]byte, error) {
	ws := writeOnlyBitStream.NewWriteOnlyBitStream()

	header := ubonHeader.NewDefaultUbonHeader(false)

	ws.AppendBitCode(header.BitCode())
	err := intEncoderHelper.WriteAutoIntToBitStream(input, ws)
	if err != nil {
		return nil, err
	}
	return ws.Bytes(), nil
}

func encodeUnsignedIntegerFirst(input any) ([]byte, error) {
	ws := writeOnlyBitStream.NewWriteOnlyBitStream()

	header := ubonHeader.NewDefaultUbonHeader(false)

	ws.AppendBitCode(header.BitCode())
	err := uIntEncoderHelper.WriteAutoUnsignedIntToBitStream(input, ws)
	if err != nil {
		return nil, err
	}
	return ws.Bytes(), nil
}

func encodeFloatFirst(input any) ([]byte, error) {

	ws := writeOnlyBitStream.NewWriteOnlyBitStream()
	header := ubonHeader.NewDefaultUbonHeader(false)
	ws.AppendBitCode(header.BitCode())

	err := floatEncoderHelper.WriteAutoFloatToBitStream(input, ws)
	if err != nil {
		return nil, err
	}
	return ws.Bytes(), nil
}

// string only input end with eos huffman and alphabet section
func encodeStringFirst(rawInput any) ([]byte, error) {

	input := rawInput.(string)

	huffFreqMap := huffman.NewHuffmanStringFrequencyMap()
	huffFreqMap.SendString(input)
	huffFreqMap.FinishСollectingStrings()

	ws := writeOnlyBitStream.NewWriteOnlyBitStream()
	header := ubonHeader.NewDefaultUbonHeader(true)
	ws.AppendBitCode(header.BitCode())
	alphabet := huffFreqMap.GetAlphabet()
	alphabetBitcode, err := huffman.AlphabetToBitcode(alphabet)
	if err != nil {
		return nil, err
	}
	ws.AppendBitCode(*alphabetBitcode)
	tree := huffman.BuildTree(alphabet)
	if tree == nil {

		return nil, errors.New("huffman tree build failed")
	}
	codes := make(map[rune]bitcode.BitCode)
	huffman.GenerateCodes(tree, bitcode.NewZeroBitCodeWithLength(0), &codes)

	err = stringEncoderHelper.WriteEncodedStringToBitStream(input, codes, ws)
	if err != nil {
		return nil, err
	}
	return ws.Bytes(), nil
}

// single bool only header with bool optcode and signle bit for boolean
func encodeBooleanFirst(input bool) ([]byte, error) {

	header := ubonHeader.NewDefaultUbonHeader(false)
	ws := writeOnlyBitStream.NewWriteOnlyBitStream()
	ws.AppendBitCode(header.BitCode())
	booleanEncoderHelper.WriteBooleanToBitStream(input, ws)
	return ws.Bytes(), nil
}

// single null only header with special opcode
func encodeNullFirst() ([]byte, error) {
	ws := writeOnlyBitStream.NewWriteOnlyBitStream()
	header := ubonHeader.NewDefaultUbonHeader(false)
	ws.AppendBitCode(header.BitCode())
	ws.AppendBitCode(ubonOpCodes.UBON_OP_NEXT_SPECIAL)
	return ws.Bytes(), nil
}
