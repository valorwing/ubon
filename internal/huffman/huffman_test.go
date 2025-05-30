package huffman_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	bitcode "ubon/internal/bitCode"
	"ubon/internal/huffman"
	"ubon/internal/readOnlyBitStream"
	"ubon/internal/writeOnlyBitStream"
)

func TestAlphabetSerialization(t *testing.T) {

	testStrings := []string{
		"hello world",
		"hello world2",
	}

	huffFreqMap := huffman.NewHuffmanStringFrequencyMap()

	for _, v := range testStrings {
		huffFreqMap.SendString(v)
	}

	huffFreqMap.FinishСollectingStrings()

	alphabet := huffFreqMap.GetAlphabet()
	fmt.Println("Generated alphabet: ", strings.Join(alphabet, ""))

	alphabetBitcode, err := huffman.AlphabetToBitcode(alphabet)
	if err != nil {
		t.Fail()
	}
	ws := writeOnlyBitStream.NewWriteOnlyBitStream()
	ws.AppendBitCode(*alphabetBitcode)
	serializedData := ws.Bytes()
	fmt.Println("Serialized alphabet bytes len: ", len(serializedData))
	rs := readOnlyBitStream.NewReadOnlyBitStream(serializedData)
	restoredAlphabet, err := huffman.AlphabetFromBitStream(&rs)
	if err != nil {
		t.Fail()
	}
	if !slices.Equal(alphabet, restoredAlphabet) {
		t.Fail()
	}
}

func TestEqualityEOSConstVarBitcodes(t *testing.T) {
	aBytes := []byte(huffman.EOS_Char)
	a := bitcode.NewBitCodeFromBytes(aBytes...)
	b := huffman.EOS_Char_Bitcode
	if !a.Equal(b) {
		t.Fail()
	}
	if a.String() != b.String() {
		t.Fail()
	}
}

func TestBaseFunctional(t *testing.T) {

	testStrings := []string{
		"hello world",
		"hello world2",
	}

	huffFreqMap := huffman.NewHuffmanStringFrequencyMap()

	for _, v := range testStrings {
		huffFreqMap.SendString(v)
	}

	huffFreqMap.FinishСollectingStrings()

	alphabet := huffFreqMap.GetAlphabet()

	fmt.Println("Alphabet: ")
	for _, v := range alphabet {
		fmt.Print(v)
	}
	fmt.Println()
	tree := huffman.BuildTree(alphabet)

	if tree == nil {

		return
	}

	codes := make(map[string]bitcode.BitCode)
	huffman.GenerateCodes(tree, bitcode.NewZeroBitCodeWithLength(0), &codes)

	fmt.Println("Codes: ")
	for k, v := range codes {
		fmt.Printf("char: %s code: %s \n", k, v.String())
	}

	bs := writeOnlyBitStream.NewWriteOnlyBitStream()
	//endcode strings sep EOS
	for _, v := range testStrings {
		for _, runeVal := range v {
			bs.AppendBitCode(codes[string([]rune{runeVal})])
		}
		bs.AppendBitCode(codes[huffman.EOS_Char])
	}

	bsBystes := bs.Bytes()
	rawLen := len([]byte(strings.Join(testStrings, huffman.EOS_Char)))
	encLen := len(bsBystes)
	fmt.Println("Raw bytes len:", rawLen)
	fmt.Println("Encoded bytes len: ", encLen)
	fmt.Println("Rate: ", float64(encLen)/float64(rawLen))
	rbs := readOnlyBitStream.NewReadOnlyBitStream(bsBystes)

	outString := make([]string, 0)
	//for test
	minBitCodeLen := -1
	reverseCodes := map[string]string{}
	for k, v := range codes {
		if minBitCodeLen == -1 {
			minBitCodeLen = v.BitLength()
		} else if minBitCodeLen > v.BitLength() {
			minBitCodeLen = v.BitLength()
		}
		reverseCodes[v.String()] = k
	}

	//check stream self control for this test allowed know size of strings
	var currentBitCode bitcode.BitCode = bitcode.NewZeroBitCodeWithLength(0)
	var currentString = ""
	isNewChar := true
	for {
		readedLen := 1
		if isNewChar {
			readedLen = minBitCodeLen
		}
		b, err := rbs.ReadBitCode(readedLen)
		if err != nil {
			t.FailNow()
			return
		}
		currentBitCode.AppendBitCode(*b)
		str, ok := reverseCodes[currentBitCode.String()]
		if !ok {
			isNewChar = false
			continue
		}
		currentBitCode.Clear()
		isNewChar = true
		if str == huffman.EOS_Char {

			outString = append(outString, currentString)
			currentString = ""
			//note!!! check stream self control for this test allowed know size of strings
			if len(outString) == len(testStrings) {
				break
			}
		} else {
			currentString = currentString + str
		}
	}

	if !slices.Equal(outString, testStrings) {
		t.Fail()
	} else {
		fmt.Println("HELL YEAH")
	}
}
