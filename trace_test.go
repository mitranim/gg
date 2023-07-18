package gg_test

import (
	"testing"

	"github.com/mitranim/gg"
)

func BenchmarkCaptureTrace_shallow(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.CaptureTrace(0))
	}
}

func BenchmarkCaptureTrace_deep(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(trace0())
	}
}

func BenchmarkTrace_Frames_shallow(b *testing.B) {
	trace := gg.CaptureTrace(0)
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(trace.Frames())
	}
}

func BenchmarkTrace_Frames_deep(b *testing.B) {
	trace := trace0()
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(trace.Frames())
	}
}

func BenchmarkFrames_NameWidth(b *testing.B) {
	frames := trace0().Frames()
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(frames.NameWidth())
	}
}

func BenchmarkFrames_AppendIndentTableTo(b *testing.B) {
	frames := trace0().Frames()
	buf := make([]byte, 0, 1<<16)
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(frames.AppendIndentTableTo(buf, 0))
	}
}

func BenchmarkFrames_AppendIndentTableTo_rel_path(b *testing.B) {
	defer gg.SnapSwap(&gg.TraceSkipLang, true).Done()
	defer gg.SnapSwap(&gg.TraceBaseDir, gg.Cwd()).Done()

	frames := trace0().Frames()
	buf := make([]byte, 0, 1<<16)
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(frames.AppendIndentTableTo(buf, 0))
	}
}

func BenchmarkTrace_capture_append(b *testing.B) {
	buf := make([]byte, 0, 1<<16)
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(trace0().AppendTo(buf))
	}
}

func trace0() gg.Trace { return trace1() }
func trace1() gg.Trace { return trace2() }
func trace2() gg.Trace { return trace3() }
func trace3() gg.Trace { return trace4() }
func trace4() gg.Trace { return trace5() }
func trace5() gg.Trace { return trace6() }
func trace6() gg.Trace { return trace7() }
func trace7() gg.Trace { return trace8() }
func trace8() gg.Trace { return trace9() }
func trace9() gg.Trace { return gg.CaptureTrace(0) }
