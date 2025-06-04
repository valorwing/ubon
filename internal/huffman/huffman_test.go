package huffman_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"
	bitcode "ubon/internal/bitCode"
	bitcodeHashMap "ubon/internal/bitCodeHashMap"
	"ubon/internal/huffman"
	"ubon/internal/readOnlyBitStream"
	"ubon/internal/writeOnlyBitStream"
)

func TestBundledOverhead(t *testing.T) {
	startTestTime := time.Now()
	testStrings := []string{
		`Once upon a time, a small fox named Felix found a forgotten flute in the forest. Curious, he played a note. To his surprise, the trees shimmered, the wind paused, and birds began to sing in harmony. "Magic!" he whispered. Every day after, the forest danced to Felix's tune.
`,
		`fox fox fox fox runs fast fast fast through the forest forest forest full of foxes and fast things
`,
	}

	huffFreqMap := huffman.NewHuffmanStringFrequencyMap()

	for _, v := range testStrings {
		huffFreqMap.SendString(v)
	}

	huffFreqMap.FinishСollectingStrings()

	alphabet := huffFreqMap.GetAlphabet()

	alphabetBitcode, err := huffman.AlphabetToBitcode(alphabet)
	if err != nil {
		t.Fail()
	}
	ws := writeOnlyBitStream.NewWriteOnlyBitStream()
	ws.AppendBitCode(*alphabetBitcode)
	serializedAplhabetData := ws.Bytes()
	fmt.Println("Serialized alphabet bytes len: ", len(serializedAplhabetData))
	rs := readOnlyBitStream.NewReadOnlyBitStream(serializedAplhabetData)
	restoredAlphabet, err := huffman.AlphabetFromBitStream(&rs)
	if err != nil {
		t.Fail()
	}
	if !slices.Equal(alphabet, restoredAlphabet) {
		t.Fail()
	}

	fmt.Println("Alphabet: ")
	for _, v := range alphabet {
		fmt.Print(v)
	}
	fmt.Println()
	tree := huffman.BuildTree(alphabet)

	if tree == nil {

		return
	}

	codes := make(map[rune]bitcode.BitCode)
	huffman.GenerateCodes(tree, bitcode.NewZeroBitCodeWithLength(0), &codes)

	fmt.Println("Codes: ")
	for k, v := range codes {
		fmt.Printf("char: %s code: %s \n", string(k), v.String())
	}

	bs := writeOnlyBitStream.NewWriteOnlyBitStream()
	//endcode strings sep EOS
	for _, v := range testStrings {
		for _, runeVal := range v {
			bs.AppendBitCode(codes[runeVal])
		}
		bs.AppendBitCode(codes[huffman.EOS_Char])
	}

	bsBystes := bs.Bytes()
	rawLen := len([]byte(strings.Join(testStrings, string([]rune{huffman.EOS_Char}))))
	encLen := len(bsBystes)
	fmt.Println("Raw bytes len:", rawLen)
	fmt.Println("Encoded bytes len: ", encLen)
	fmt.Println("Rate: ", float64(encLen)/float64(rawLen))
	rbs := readOnlyBitStream.NewReadOnlyBitStream(bsBystes)

	outString := make([]string, 0)
	//for test
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

	//check stream self control for this test allowed know size of strings
	var currentBitCode bitcode.BitCode = bitcode.NewZeroBitCodeWithLength(0)
	var currentString = []rune{}
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
		runeVal, ok := reverseCodes.Get(currentBitCode)
		if !ok {
			isNewChar = false
			continue
		}
		currentBitCode.Clear()
		isNewChar = true
		if runeVal == huffman.EOS_Char {

			outString = append(outString, string(currentString))
			currentString = currentString[:0]
			//note!!! check stream self control for this test allowed know size of strings
			if len(outString) == len(testStrings) {
				break
			}
		} else {
			currentString = append(currentString, runeVal)
		}
	}

	if !slices.Equal(outString, testStrings) {
		t.Fail()
	}

	//compare bundled encoded huffman alphabe and payload

	fullHuffmanData := append(serializedAplhabetData, bsBystes...)
	testData := []byte(strings.Join(testStrings, ""))
	fmt.Println("Compress rate with overhead: ", float64(len(fullHuffmanData))/float64(len(testData)))
	fmt.Println("Test execution time: ", time.Since(startTestTime))
}
