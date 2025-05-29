package huffman

type HuffmanHeap[T HuffmanCharInterface] []*HuffmanNode[T]

func (h HuffmanHeap[T]) Len() int {
	return len(h)
}
func (h HuffmanHeap[T]) Less(i, j int) bool {
	return h[i].frequency < h[j].frequency
}
func (h HuffmanHeap[T]) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *HuffmanHeap[T]) Push(x interface{}) {
	*h = append(*h, x.(*HuffmanNode[T]))
}

func (h *HuffmanHeap[T]) Pop() interface{} {
	var retVal interface{}
	hLen := len(*h)
	retVal, *h = (*h)[hLen-1], (*h)[:hLen-1]
	return retVal
}
