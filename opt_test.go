package gg_test

import (
	"testing"

	"github.com/mitranim/gg"
)

func BenchmarkOpt_String(b *testing.B) {
	val := gg.OptVal(`str`)

	for i := 0; i < b.N; i++ {
		gg.Nop1(val.String())
	}
}
