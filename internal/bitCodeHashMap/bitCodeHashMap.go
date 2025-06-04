package bitcodeHashMap

import (
	bitcode "ubon/internal/bitCode"
)

type entry[V any] struct {
	hashKey uint64
	key     bitcode.BitCode
	value   V
}

type BitcodeHashMap[V any] struct {
	buckets     [][]entry[V]
	bucketCount int
	size        int
}

// bucket count 0.25/0.5 final length
func NewHashMap[V any](bucketCount int) *BitcodeHashMap[V] {
	if bucketCount <= 0 {
		bucketCount = 16
	}
	return &BitcodeHashMap[V]{
		buckets:     make([][]entry[V], bucketCount),
		bucketCount: bucketCount,
		size:        0,
	}
}

func (h *BitcodeHashMap[V]) Put(key bitcode.BitCode, value V) {

	hk := key.Hash()
	idx := int(hk % uint64(h.bucketCount))
	bucket := h.buckets[idx]

	for i, ent := range bucket {
		if ent.hashKey == hk && ent.key.Equal(key) {
			h.buckets[idx][i].value = value
			return
		}
	}

	h.buckets[idx] = append(bucket, entry[V]{
		hashKey: hk,
		key:     key,
		value:   value,
	})
	h.size++
}

func (h *BitcodeHashMap[V]) Get(key bitcode.BitCode) (V, bool) {
	hk := key.Hash()
	idx := int(hk % uint64(h.bucketCount))
	bucket := h.buckets[idx]

	for _, ent := range bucket {
		if ent.hashKey == hk && ent.key.Equal(key) {
			return ent.value, true
		}
	}
	var zero V
	return zero, false
}

func (h *BitcodeHashMap[V]) Delete(key bitcode.BitCode) bool {
	hk := key.Hash()
	idx := int(hk % uint64(h.bucketCount))
	bucket := h.buckets[idx]

	for i, ent := range bucket {
		if ent.hashKey == hk && ent.key.Equal(key) {
			last := len(bucket) - 1
			bucket[i] = bucket[last]
			h.buckets[idx] = bucket[:last]
			h.size--
			return true
		}
	}
	return false
}

func (h *BitcodeHashMap[V]) Len() int {

	return h.size
}

func (h *BitcodeHashMap[V]) Keys() []bitcode.BitCode {

	keys := make([]bitcode.BitCode, 0, h.size)

	for _, bucket := range h.buckets {
		for _, ent := range bucket {
			keys = append(keys, ent.key)
		}
	}
	return keys
}
