package huffman_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	bitcode "ubon/internal/bitCode"
	"ubon/internal/huffman"
)

type MyChar string

func TestBaseFunctional(t *testing.T) {

	testString := "hello world"
	fmt.Println("input: ", testString)
	freqMap := map[MyChar]int{}

	for _, runeVal := range testString {
		count, ok := freqMap[MyChar(string([]rune{runeVal}))]
		if !ok {
			freqMap[MyChar(string([]rune{runeVal}))] = 1
		} else {
			freqMap[MyChar(string([]rune{runeVal}))] = count + 1
		}
	}

	alphabet := []MyChar{}

	for k := range freqMap {
		alphabet = append(alphabet, MyChar(k))
	}
	slices.SortFunc(alphabet, func(a, b MyChar) int {
		if freqMap[a] < freqMap[b] {
			return -1
		} else if freqMap[a] > freqMap[b] {
			return 1
		} else {
			return strings.Compare(string(a), string(b))
		}
	})

	tree := huffman.BuildTree(alphabet)

	if tree == nil {

		return
	}

	codes := make(map[MyChar]bitcode.BitCode)
	huffman.GenerateCodes(tree, bitcode.NewZeroBitCodeWithLength(0), &codes)
	fmt.Println("Alphabet: ")
	for _, v := range alphabet {
		fmt.Print(v)
	}
	fmt.Println()
	fmt.Println("Codes: ")
	for k, v := range codes {
		fmt.Printf("char: %s code: %s \n", k, v.String())
	}

}
