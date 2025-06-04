package bitcodeHashMap_test

import (
	"fmt"
	"testing"
	"time"
	bitcode "ubon/internal/bitCode"
	bitcodeHashMap "ubon/internal/bitCodeHashMap"
	"ubon/internal/huffman"
)

func TestIntegrationHuffman(t *testing.T) {
	testStrings := []string{
		"hello world",
		"hello world2",
		"oalalalala",
		"uffff",
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

	codes := make(map[rune]bitcode.BitCode)
	huffman.GenerateCodes(tree, bitcode.NewZeroBitCodeWithLength(0), &codes)

	reverseCodes := bitcodeHashMap.NewHashMap[rune](len(codes) / 2)

	for k, v := range codes {
		reverseCodes.Put(v, k)
	}

	//test
	for k, v := range codes {
		if k1, ok := reverseCodes.Get(v); ok {
			if k != k1 {
				t.Fail()
			}
		} else {
			t.Fail()
		}
	}

	//compare access speed
	reverseStart := time.Now()
	for _, k := range reverseCodes.Keys() {
		_, _ = reverseCodes.Get(k)
	}
	reverseTime := time.Since(reverseStart)

	codesStart := time.Now()
	for k := range codes {
		_ = codes[k]
	}
	codesTime := time.Since(codesStart)

	fmt.Println("BitcodeHashMap access time: ", reverseTime)
	fmt.Println("Buildin map access time: ", codesTime)
	fmt.Println("Diff: ", float64(reverseTime.Nanoseconds())/float64(codesTime.Nanoseconds()))
}
