package bitcode

import (
	"encoding/binary"
)

const (
	GridLength      = 4
	CellCapacity    = 64                        // 64 bit cell capacity
	GridBitCapacity = GridLength * CellCapacity // 4 * 64 Grid
	CellMaxIndex    = CellCapacity - 1          // 63 bits for edit, last bit is for next cell
)

type BitCodeStorageGridCell struct {
	CellGrid     [GridLength]uint64
	NextCellGrid *BitCodeStorageGridCell
}

type BitCode struct {
	headCellGrid         *BitCodeStorageGridCell
	tailCellGrid         *BitCodeStorageGridCell
	tailCellGridBitIndex uint8
}

func NewZeroBitCodeWithLength(length int) BitCode {

	if length == 0 {
		cellGrid := &BitCodeStorageGridCell{
			CellGrid:     [4]uint64{0, 0, 0, 0},
			NextCellGrid: nil,
		}
		return BitCode{
			headCellGrid:         cellGrid,
			tailCellGrid:         cellGrid,
			tailCellGridBitIndex: 0,
		}
	} else if length < GridBitCapacity {
		cellGrid := &BitCodeStorageGridCell{
			CellGrid:     [4]uint64{0, 0, 0, 0},
			NextCellGrid: nil,
		}
		return BitCode{
			headCellGrid:         cellGrid,
			tailCellGrid:         cellGrid,
			tailCellGridBitIndex: uint8(length),
		}
	} else {

		cellGridCount := length / GridBitCapacity

		tailCellGrid := &BitCodeStorageGridCell{
			CellGrid:     [4]uint64{0, 0, 0, 0},
			NextCellGrid: nil,
		}
		headCellGrid := tailCellGrid
		//Always + 1 alloc
		for i := 0; i < cellGridCount; i++ {
			newCellGrid := &BitCodeStorageGridCell{
				CellGrid:     [4]uint64{0, 0, 0, 0},
				NextCellGrid: nil,
			}
			tailCellGrid.NextCellGrid = newCellGrid
			tailCellGrid = newCellGrid
		}
		return BitCode{
			headCellGrid:         headCellGrid,
			tailCellGrid:         tailCellGrid,
			tailCellGridBitIndex: uint8(length % GridBitCapacity),
		}

	}
}

func NewBitCodeWithBoolCode(code ...bool) BitCode {
	retVal := NewZeroBitCodeWithLength(len(code))
	for i, v := range code {
		if v {
			jumpCount := i / GridBitCapacity
			cellToEdit := retVal.headCellGrid
			for j := 0; j < jumpCount; j++ {
				cellToEdit = cellToEdit.NextCellGrid
			}
			cellGridBitIndex := i % GridBitCapacity
			cellGridIndex := cellGridBitIndex / CellCapacity
			cellBitIndex := cellGridBitIndex % CellCapacity

			cellToEdit.CellGrid[cellGridIndex] |= 1 << (CellMaxIndex - cellBitIndex)
		}
	}
	return retVal
}

func NewBitCodeFromBytes(bytesData ...byte) BitCode {
	retVal := NewZeroBitCodeWithLength(0)
	bitPoiner := 0
	for _, v := range bytesData {
		for i := 0; i < 8; i++ {

			retVal.Append((v & (1 << (7 - uint(i%8)))) != 0)

			bitPoiner++
		}
	}
	return retVal
}

func (b BitCode) BitLength() int {

	if b.headCellGrid == nil {
		return 0
	} else if b.headCellGrid.NextCellGrid == nil {
		return int(b.tailCellGridBitIndex)
	} else {
		retVal := GridBitCapacity
		currentGrid := b.headCellGrid.NextCellGrid
		for {
			if currentGrid.NextCellGrid != nil {
				retVal += GridBitCapacity
				currentGrid = currentGrid.NextCellGrid
			} else {
				retVal += int(b.tailCellGridBitIndex + 1)
				break
			}
		}
		return retVal
	}
}

func (b *BitCode) SetBit(bitIndex int) {

	if bitIndex >= b.BitLength() {
		panic("BitCode: trying to set bit out of range")
	}
	jumpCount := bitIndex / GridBitCapacity
	cellToEdit := b.headCellGrid
	for j := 0; j < jumpCount; j++ {
		cellToEdit = cellToEdit.NextCellGrid
	}
	if cellToEdit == nil {
		panic("BitCode: trying to set bit in non-existing cell")
	}
	cellGridBitIndex := bitIndex % GridBitCapacity
	cellGridIndex := cellGridBitIndex / CellCapacity
	cellBitIndex := cellGridBitIndex % CellCapacity

	cellToEdit.CellGrid[cellGridIndex] |= 1 << (CellMaxIndex - cellBitIndex)
}

func (b BitCode) Clone() BitCode {

	if b.headCellGrid == nil {
		return BitCode{headCellGrid: nil, tailCellGrid: nil, tailCellGridBitIndex: 0}
	} else if b.headCellGrid.NextCellGrid == nil {
		newHeadCellGrid := &BitCodeStorageGridCell{
			CellGrid:     b.headCellGrid.CellGrid,
			NextCellGrid: nil,
		}
		return BitCode{headCellGrid: newHeadCellGrid, tailCellGrid: newHeadCellGrid, tailCellGridBitIndex: b.tailCellGridBitIndex}
	} else {

		newHeadCellGrid := &BitCodeStorageGridCell{
			CellGrid:     b.headCellGrid.CellGrid,
			NextCellGrid: nil,
		}
		newTailCellGrid := newHeadCellGrid

		currentReadedCell := b.headCellGrid.NextCellGrid
		for {

			internalNewCellGrid := &BitCodeStorageGridCell{
				CellGrid:     currentReadedCell.CellGrid,
				NextCellGrid: nil,
			}
			newTailCellGrid.NextCellGrid = internalNewCellGrid
			newTailCellGrid = internalNewCellGrid

			if currentReadedCell.NextCellGrid == nil {
				break
			} else {
				currentReadedCell = currentReadedCell.NextCellGrid
			}
		}
		return BitCode{headCellGrid: newHeadCellGrid, tailCellGrid: newTailCellGrid, tailCellGridBitIndex: b.tailCellGridBitIndex}
	}
}

func (b BitCode) CloneAppend(bit bool) BitCode {
	retVal := b.Clone()
	retVal.Append(bit)
	return retVal
}

func (b *BitCode) AppendBitCode(other BitCode) {
	//TODO: faster append
	for i := 0; i < other.BitLength(); i++ {
		b.Append(other.GetBit(i))
	}
}

func (b *BitCode) Append(bit bool) {
	if b.tailCellGrid == nil {
		b.headCellGrid = &BitCodeStorageGridCell{
			CellGrid: [GridLength]uint64{0, 0, 0, 0},
		}
		b.tailCellGrid = b.headCellGrid
		b.tailCellGridBitIndex = 0
	}

	if bit {
		b.tailCellGrid.CellGrid[b.tailCellGridBitIndex/CellCapacity] |= 1 << (CellMaxIndex - b.tailCellGridBitIndex%CellCapacity)
	}
	b.tailCellGridBitIndex++
	if b.tailCellGridBitIndex == GridBitCapacity-1 {
		// If the current cell is full, we need to create a new one
		newCellGrid := &BitCodeStorageGridCell{
			CellGrid: [GridLength]uint64{0, 0, 0, 0},
		}

		b.tailCellGrid.NextCellGrid = newCellGrid
		b.tailCellGrid = newCellGrid
		b.tailCellGridBitIndex = 0
	}
}

func (b BitCode) GetBit(i int) bool {
	if i < 0 || i >= b.BitLength() {
		panic("BitCode: trying to get bit out of range")
	}
	jumpCount := i / GridBitCapacity
	gridBitIndex := i % GridBitCapacity
	gridIndex := gridBitIndex / CellCapacity

	cellToRead := b.headCellGrid
	for j := 0; j < jumpCount; j++ {
		if cellToRead.NextCellGrid == nil {
			panic("BitCode: trying to get bit from non-existing cell")
		}
		cellToRead = cellToRead.NextCellGrid
	}
	return (cellToRead.CellGrid[gridIndex] & (1 << (CellMaxIndex - gridBitIndex%CellCapacity))) != 0
}

// debug
func (b BitCode) String() string {
	//no lib
	retVal := ""
	for i := 0; i < b.BitLength(); i++ {
		s := "0"
		if b.GetBit(i) {
			s = "1"
		}
		retVal = retVal + s
	}
	return retVal
}

func (b *BitCode) Clear() {
	b.headCellGrid = nil
	b.tailCellGrid = nil
	b.tailCellGridBitIndex = 0
}

func (b BitCode) Equal(other BitCode) bool {
	//TODO: optimize
	if b.headCellGrid == nil && other.headCellGrid == nil {
		return true
	} else if b.BitLength() != other.BitLength() {
		return false
	} else {
		for i := 0; i < b.BitLength(); i++ {
			if b.GetBit(i) != other.GetBit(i) {
				return false
			}
		}
	}
	return true
}

func (b BitCode) Bytes() []byte {
	if b.headCellGrid == nil {
		return nil
	}
	retVal := make([]byte, 0, b.BitLength()/8+8)

	cell := b.headCellGrid
	for cell != nil {
		buf := make([]byte, GridLength*8)
		for i, v := range cell.CellGrid {
			binary.BigEndian.PutUint64(buf[i*8:], v)
		}
		retVal = append(retVal, buf...)
		cell = cell.NextCellGrid
	}
	//adjust length
	finalLength := b.BitLength() / 8
	if b.BitLength()%8 > 0 {
		finalLength++
	}
	retVal = retVal[:finalLength]
	return retVal
}

// !!!Warning for cell grid 4 x 64 for other - refactor
func (b BitCode) Hash() uint64 {
	var h uint64 = 14695981039346656037
	const prime uint64 = 1099511628211

	cell := b.headCellGrid
	for cell != nil {
		for i, cellValue := range cell.CellGrid {
			if cell.NextCellGrid != nil || int(b.tailCellGridBitIndex) > i*CellCapacity {
				h ^= cellValue
				h *= prime
			}
		}
		cell = cell.NextCellGrid
	}

	h ^= uint64(b.BitLength())
	h *= prime
	return h
}
