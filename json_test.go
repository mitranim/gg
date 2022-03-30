package gg_test

import (
	"encoding/json"
	"testing"

	"github.com/mitranim/gg"
)

func Benchmark_json_Marshal(b *testing.B) {
	var val SomeModel

	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.Try1(json.Marshal(val)))
	}
}

func BenchmarkJsonBytes(b *testing.B) {
	var val SomeModel

	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.JsonBytes(val))
	}
}

func Benchmark_json_Marshal_string(b *testing.B) {
	var val SomeModel

	for i := 0; i < b.N; i++ {
		gg.Nop1(string(gg.Try1(json.Marshal(val))))
	}
}

func BenchmarkJsonString(b *testing.B) {
	var val SomeModel

	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.JsonString(val))
	}
}

func Benchmark_json_Unmarshal(b *testing.B) {
	var val int

	for i := 0; i < b.N; i++ {
		gg.Try(json.Unmarshal(gg.ToBytes(`123`), &val))
	}
}

func BenchmarkJsonParseTo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.JsonParseTo[int](`123`))
	}
}

func BenchmarkJsonParse(b *testing.B) {
	var val int

	for i := 0; i < b.N; i++ {
		gg.JsonParse(`123`, &val)
	}
}
