package gg_test

import (
	"encoding/json"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func Benchmark_json_Marshal(b *testing.B) {
	var val SomeModel

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Try1(json.Marshal(val)))
	}
}

func BenchmarkJsonBytes(b *testing.B) {
	var val SomeModel

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.JsonBytes(val))
	}
}

func Benchmark_json_Marshal_string(b *testing.B) {
	var val SomeModel

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(string(gg.Try1(json.Marshal(val))))
	}
}

func BenchmarkJsonString(b *testing.B) {
	var val SomeModel

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.JsonString(val))
	}
}

func Benchmark_json_Unmarshal(b *testing.B) {
	var val int

	for ind := 0; ind < b.N; ind++ {
		gg.Try(json.Unmarshal(gg.ToBytes(`123`), &val))
	}
}

func BenchmarkJsonDecodeTo(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.JsonDecodeTo[int](`123`))
	}
}

func BenchmarkJsonDecode(b *testing.B) {
	var val int

	for ind := 0; ind < b.N; ind++ {
		gg.JsonDecode(`123`, &val)
	}
}

func TestJsonDecodeTo(t *testing.T) {
	gtest.Catch(t)

	gtest.Eq(
		gg.JsonDecodeTo[SomeModel](`{"id":10}`),
		SomeModel{Id: 10},
	)

	gtest.Eq(
		gg.JsonDecodeTo[SomeModel]([]byte(`{"id":10}`)),
		SomeModel{Id: 10},
	)
}
