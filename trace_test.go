package gg_test

import (
	"strings"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

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

func TestTrace_capture_inlined_frames(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Str(trace0(), `
gg_test.trace9                           trace_test.go:20
gg_test.trace8                           trace_test.go:19
gg_test.trace7                           trace_test.go:18
gg_test.trace6                           trace_test.go:17
gg_test.trace5                           trace_test.go:16
gg_test.trace4                           trace_test.go:15
gg_test.trace3                           trace_test.go:14
gg_test.trace2                           trace_test.go:13
gg_test.trace1                           trace_test.go:12
gg_test.trace0                           trace_test.go:11
gg_test.TestTrace_capture_inlined_frames trace_test.go:25`)
}

func TestTrace_small_buffer_long_trace(t *testing.T) {
	defer gtest.Catch(t)

	t.Cleanup(gg.SnapSwap(&gg.TraceBufSize, 4).Done)

	strings.Contains(trace0().String(), `gg_test.trace9`)
}

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
