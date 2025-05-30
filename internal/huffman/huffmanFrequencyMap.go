package huffman

import (
	"context"
	"slices"
	"strings"
	"sync"
)

const (
	receiveBufLen int = 300
)

type HuffmanStringFrequencyMap struct {
	frequencyMap              map[string]uint32
	receiveStringStatsChannel chan string
	ctx                       context.Context
	ctxCancelFunction         func()
	alphabet                  []string
	alphabetCondMutex         sync.Mutex
	alphabetCond              sync.Cond
	alphabetReady             bool
}

func NewHuffmanStringFrequencyMap() *HuffmanStringFrequencyMap {
	ctx, ctxCancelFunction := context.WithCancel(context.Background())

	retVal := &HuffmanStringFrequencyMap{
		frequencyMap:              map[string]uint32{},
		receiveStringStatsChannel: make(chan string, receiveBufLen),
		ctx:                       ctx,
		ctxCancelFunction:         ctxCancelFunction,
		alphabetCondMutex:         sync.Mutex{},
		alphabetReady:             false,
	}
	retVal.alphabetCond = *sync.NewCond(&retVal.alphabetCondMutex)
	go retVal.collectLoop()
	return retVal
}

func (fm *HuffmanStringFrequencyMap) SendString(value string) {
	fm.receiveStringStatsChannel <- value
}

func (fm *HuffmanStringFrequencyMap) FinishСollectingStrings() {
	fm.ctxCancelFunction()
}

func (fm *HuffmanStringFrequencyMap) GetAlphabet() []string {
	fm.alphabetCond.L.Lock()
	defer fm.alphabetCond.L.Unlock()
	for !fm.alphabetReady {
		fm.alphabetCond.Wait()
	}
	return fm.alphabet
}

func (fm *HuffmanStringFrequencyMap) collectLoop() {
	for {
		select {
		case <-fm.ctx.Done():
			close(fm.receiveStringStatsChannel)
			for v := range fm.receiveStringStatsChannel {
				fm.collectValue(v)
			}

			fm.finishAndBuildAlphabet()
			return

		case val := <-fm.receiveStringStatsChannel:
			fm.collectValue(val)
		}
	}
}

func (fm *HuffmanStringFrequencyMap) collectValue(value string) {

	for _, runeVal := range value {
		fm.frequencyMap[string([]rune{runeVal})]++
	}
	//Add EOS to stat after any string
	fm.frequencyMap[EOS_Char]++
}

func (fm *HuffmanStringFrequencyMap) finishAndBuildAlphabet() {
	alphabet := []string{}
	for k := range fm.frequencyMap {
		alphabet = append(alphabet, k)
	}
	slices.SortFunc(alphabet, func(a, b string) int {
		if fm.frequencyMap[a] < fm.frequencyMap[b] {
			return -1
		} else if fm.frequencyMap[a] > fm.frequencyMap[b] {
			return 1
		} else {
			return strings.Compare(string(a), string(b))
		}
	})
	fm.alphabetCond.L.Lock()
	defer fm.alphabetCond.L.Unlock()
	fm.alphabet = alphabet
	fm.alphabetReady = true
	fm.alphabetCond.Broadcast()

}
