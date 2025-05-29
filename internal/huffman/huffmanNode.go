package huffman

type HuffmanNode[T HuffmanCharInterface] struct {
	Char      T
	frequency uint16
	Left      *HuffmanNode[T]
	Right     *HuffmanNode[T]
}
