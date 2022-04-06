package gg

import (
	rt "runtime"
	"strings"
)

// These variables control how stack traces are printed.
var (
	TraceTable     = true
	TraceSkipLang  = true
	TraceShortName = true
	TraceBaseDir   = `` // Set to `Cwd()` for better traces.
)

// Free cast of the given ptr slice to `Trace`.
func ToTrace(val []uintptr) Trace { return CastUnsafe[Trace](val) }

/*
Shortcut for capturing a trace of the the current call stack, skipping N frames
where 1 corresponds to the caller's frame.
*/
func CaptureTrace(skip int) Trace { return make(Trace, 64).Capture(skip + 1) }

/*
Alias of `[]uintptr` with various methods for capturing and printing
stack traces.
*/
type Trace []Caller

/*
Uses `runtime.Callers` to capture the current call stack into the given `Trace`,
which must have enough capacity. The returned slice is truncated.
*/
func (self Trace) Capture(skip int) Trace {
	return self[:rt.Callers(skip+2, self.Prim())]
}

/*
Returns a multi-line text representation of the trace, with no leading
indentation. See `.AppendIndent`.
*/
func (self Trace) String() string { return AppenderString(self) }

/*
Appends a multi-line text representation of the trace, with no leading
indentation. See `.AppendIndent`.
*/
func (self Trace) Append(buf []byte) []byte { return self.AppendIndent(buf, 0) }

/*
Returns a multi-line text representation of the trace with the given leading
indentation. See `.AppendIndent`.
*/
func (self Trace) StringIndent(lvl int) string {
	return ToString(self.AppendIndent(nil, lvl))
}

/*
Appends a multi-line text representation of the trace, with the given leading
indentation. Used internally by other trace printing methods. Affected by the
various "Trace*" variables. If `TraceTable` is true, the trace is formatted as
a table, where each frame takes only one line, and names are aligned.
Otherwise, the trace is formatted similarly to the default representation used
by the Go runtime.
*/
func (self Trace) AppendIndent(buf []byte, lvl int) []byte {
	if TraceTable {
		return self.AppendIndentTable(buf, lvl)
	}
	return self.AppendIndentMulti(buf, lvl)
}

/*
Appends a table-style representation of the trace. Used internally by
`.AppendIndent` if `TraceTable` is true.
*/
func (self Trace) AppendIndentTable(buf []byte, lvl int) []byte {
	return self.Frames().AppendIndentTable(buf, lvl)
}

/*
Appends a representation of the trace similar to the default used by the Go
runtime. Used internally by `.AppendIndent` if `TraceTable` is false.
*/
func (self Trace) AppendIndentMulti(buf []byte, lvl int) []byte {
	for _, val := range self {
		buf = val.AppendNewlineIndent(buf, lvl)
	}
	return buf
}

/*
Returns a table-style representation of the trace with the given leading
indentation.
*/
func (self Trace) TableIndent(lvl int) string {
	return ToString(self.AppendIndentTable(nil, lvl))
}

/*
Returns a table-style representation of the trace with no leading indentation.
*/
func (self Trace) Table() string { return self.TableIndent(0) }

// True if there are any non-empty frames.
func (self Trace) HasLen() bool { return Some(self, IsNonZero[Caller]) }

// Converts to `Frames`, which is used for formatting.
func (self Trace) Frames() Frames { return Map(self, Caller.Frame) }

/*
Free cast to the underlying type. Useful for `runtime.Callers` and for
implementing `StackTraced` in error types.
*/
func (self Trace) Prim() []uintptr { return CastUnsafe[[]uintptr](self) }

// Represents an entry in a call stack. Used for formatting.
type Caller uintptr

// Short for "program counter".
func (self Caller) Pc() uintptr {
	if IsZero(self) {
		return 0
	}
	// For historic reasons.
	return uintptr(self) - 1
}

// Uses `runtime.FuncForPC` to return the function corresponding to this frame.
func (self Caller) Func() *rt.Func {
	if IsZero(self) {
		return nil
	}
	return rt.FuncForPC(self.Pc())
}

// Converts to `Frame`, which is used for formatting.
func (self Caller) Frame() (out Frame) {
	out.Init(self)
	return
}

/*
Returns a single-line representation of the frame that includes function name,
file path, and row.
*/
func (self Caller) String() string { return AppenderString(self) }

func (self Caller) Append(buf []byte) []byte {
	return self.Frame().Append(buf)
}

func (self Caller) AppendIndent(buf []byte, lvl int) []byte {
	return self.Frame().AppendIndent(buf, lvl)
}

func (self Caller) AppendNewlineIndent(buf []byte, lvl int) []byte {
	return self.Frame().AppendNewlineIndent(buf, lvl)
}

type Frames []Frame

func (self Frames) NameWidth() int {
	var out int
	for _, val := range self {
		if !val.Skip() {
			out = MaxPrim2(out, len(val.NameShort()))
		}
	}
	return out
}

func (self Frames) AppendIndentTable(buf []byte, lvl int) []byte {
	wid := self.NameWidth()
	for _, val := range self {
		buf = val.AppendRowIndent(buf, lvl, wid)
	}
	return buf
}

// Represents a stack frame. Generated by `Caller`. Used for formatting.
type Frame struct {
	Caller Caller
	Func   *rt.Func
	Name   string
	File   string
	Line   int
}

// True if the frame has a known associated function.
func (self Frame) IsValid() bool { return self.Func != nil }

func (self *Frame) Init(val Caller) {
	self.Caller = val

	fun := val.Func()
	self.Func = fun

	if fun != nil {
		self.Name = FuncNameBase(fun)
		self.File, self.Line = fun.FileLine(val.Pc())
	}
}

/*
Returns a single-line representation of the frame that includes function name,
file path, and row.
*/
func (self Frame) String() string { return AppenderString(self) }

/*
Appends a single-line representation of the frame that includes function name,
file path, and row.
*/
func (self Frame) Append(inout []byte) []byte {
	buf := Buf(inout)
	if self.Skip() {
		return buf
	}

	buf.AppendString(self.NameShort())
	buf.AppendSpace()
	buf.AppendString(self.Path())
	buf.AppendString(`:`)
	buf.AppendInt(self.Line)
	return buf
}

func (self Frame) AppendIndent(inout []byte, lvl int) []byte {
	buf := Buf(inout)
	if self.Skip() {
		return buf
	}

	buf.AppendString(self.NameShort())
	buf.AppendNewline()
	buf.AppendIndents(lvl)
	buf.AppendString(self.Path())
	buf.AppendString(`:`)
	buf.AppendInt(self.Line)
	return buf
}

func (self Frame) AppendNewlineIndent(inout []byte, lvl int) []byte {
	buf := Buf(inout)
	if self.Skip() {
		return buf
	}

	buf.AppendNewline()
	buf.AppendIndents(lvl)
	return self.AppendIndent(buf, lvl+1)
}

func (self Frame) AppendRowIndent(inout []byte, lvl, wid int) []byte {
	buf := Buf(inout)
	if self.Skip() {
		return buf
	}

	name := self.NameShort()

	buf.AppendNewline()
	buf.AppendIndents(lvl)
	buf.AppendString(name)
	buf.AppendSpace()
	buf.AppendSpaces(wid - len(name))
	buf.AppendString(self.Path())
	buf.AppendString(`:`)
	buf.AppendInt(self.Line)
	return buf
}

/*
True if the frame should not be displayed, either because it's invalid, or
because `TraceSkipLang` is set and the frame represents a "language" frame
which is mostly not useful for debugging app code.
*/
func (self *Frame) Skip() bool {
	return !self.IsValid() || (TraceSkipLang && self.IsLang())
}

/*
True if the frame represents a "language" frame which is mostly not useful for
debugging app code.
*/
func (self *Frame) IsLang() bool {
	pkg := self.Pkg()
	return pkg == `runtime` || pkg == `testing`
}

// Returns the package name of the given frame.
func (self *Frame) Pkg() string {
	name := self.Name
	ind := strings.IndexByte(name, '.')
	if ind >= 0 {
		return name[:ind]
	}
	return name
}

func (self *Frame) NameShort() string {
	if TraceShortName {
		return FuncNameShort(self.Name)
	}
	return self.Name
}

func (self *Frame) Path() string {
	if TraceBaseDir != `` {
		return relOpt(TraceBaseDir, self.File)
	}
	return self.File
}
