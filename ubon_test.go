package ubon_test

import (
	"fmt"
	"testing"
	"ubon/internal/ubonDecoder"
	"ubon/internal/ubonEncoder"
)

func TestRoundTrip(t *testing.T) {

	cases := []any{
		nil,
		true, false,
		42, int64(1 << 40),
		uint(123456),
		3.14, float64(2.718),
		"hello world",
		map[string]any{
			"hello":  "world",
			"number": 123,
			"bool":   true,
		},
		map[string]any{
			"hello":  "world",
			"number": 123,
			"bool":   true,
			"object": map[string]any{
				"hello":  "world",
				"number": 123,
				"bool":   true,
			},
		},
	}

	for _, input := range cases {
		encoded, err := ubonEncoder.Encode(input)
		if err != nil {
			t.Errorf("encode failed: %v", err)
			continue
		}
		decoded, err := ubonDecoder.Decode(encoded)
		if err != nil {
			t.Errorf("decode failed: %v", err)
			continue
		}
		if fmt.Sprintf("%v", decoded) != fmt.Sprintf("%v", input) {
			t.Errorf("roundtrip mismatch: input=%v decoded=%v", input, decoded)
		}
	}
}
