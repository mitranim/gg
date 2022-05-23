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

func BenchmarkJsonParseTo(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.JsonParseTo[int](`123`))
	}
}

func BenchmarkJsonParse(b *testing.B) {
	var val int

	for ind := 0; ind < b.N; ind++ {
		gg.JsonParse(`123`, &val)
	}
}

func TestJsonParseTo(t *testing.T) {
	gtest.Catch(t)

	gtest.Eq(
		gg.JsonParseTo[SomeModel](`{"id":"10"}`),
		SomeModel{Id: `10`},
	)

	gtest.Eq(
		gg.JsonParseTo[SomeModel]([]byte(`{"id":"10"}`)),
		SomeModel{Id: `10`},
	)
}
