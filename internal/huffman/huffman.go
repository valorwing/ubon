package huffman

import (
	"container/heap"
	bitcode "ubon/internal/bitCode"
)

func BuildTree[T HuffmanCharInterface](treeSource []T) *HuffmanNode[T] {

	h := &HuffmanHeap[T]{}

	for freq, char := range treeSource {
		h.Push(&HuffmanNode[T]{Char: char, frequency: uint16(freq + 1)})
	}
	heap.Init(h)
	if h.Len() == 0 {
		return nil
	}
	for h.Len() > 1 {
		left := heap.Pop(h).(*HuffmanNode[T])
		right := heap.Pop(h).(*HuffmanNode[T])
		heap.Push(h, &HuffmanNode[T]{
			frequency: left.frequency + right.frequency,
			Left:      left,
			Right:     right,
		})
	}
	return heap.Pop(h).(*HuffmanNode[T])
}

func GenerateCodes[T HuffmanCharInterface](node *HuffmanNode[T], prefix bitcode.BitCode, codes *map[T]bitcode.BitCode) {
	if codes == nil {
		tmp := make(map[T]bitcode.BitCode)
		codes = &tmp
	}
	if node == nil {
		return
	}

	if node.Left == nil && node.Right == nil {
		(*codes)[node.Char] = prefix
	}

	GenerateCodes(node.Left, prefix.CloneAppend(false), codes)
	GenerateCodes(node.Right, prefix.CloneAppend(true), codes)
}
