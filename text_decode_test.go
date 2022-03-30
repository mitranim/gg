package gg_test

import (
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func TestParseTo(t *testing.T) {
	gtest.Eq(gg.ParseTo[int](`123`), 123)
}

func BenchmarkParseTo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.ParseTo[int](`123`))
	}
}

func BenchmarkParse(b *testing.B) {
	var val int

	for i := 0; i < b.N; i++ {
		gg.Parse(`123`, &val)
	}
}
