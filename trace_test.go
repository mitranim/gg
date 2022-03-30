package gg_test

import (
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func BenchmarkCaptureTrace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.CaptureTrace(0))
	}
}

func BenchmarkTraceAppend(b *testing.B) {
	gtest.Swap(b, &gg.TraceRelPath, false)

	trace := gg.CaptureTrace(0)
	buf := make([]byte, 0, 4096)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		gg.Nop1(trace.Append(buf))
	}
}

func BenchmarkCaptureTrace_Append(b *testing.B) {
	gtest.Swap(b, &gg.TraceRelPath, false)

	buf := make([]byte, 0, 4096)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.CaptureTrace(0).Append(buf))
	}
}
