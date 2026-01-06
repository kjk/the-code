package main

import (
	"encoding/json"
	"testing"

	"github.com/alpkeskin/gotoon"
	"github.com/toon-format/toon-go"
)

func BenchmarkJSONMarshalCompact(b *testing.B) {
	for b.Loop() {
		_, err := json.Marshal(testConfig)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONMarshal(b *testing.B) {
	for b.Loop() {
		_, err := json.MarshalIndent(testConfig, "", "  ")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkToonMarshal(b *testing.B) {

	for b.Loop() {
		_, _err := toon.Marshal(testConfig)
		if _err != nil {
			b.Fatal(_err)
		}
	}
}

func BenchmarkGotoonEncode(b *testing.B) {
	for b.Loop() {
		_, _err := gotoon.Encode(testConfig)
		if _err != nil {
			b.Fatal(_err)
		}
	}
}
