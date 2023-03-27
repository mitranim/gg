package gg_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func TestGraphDir_invalid_missing_deps(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(
		`unable to build dependency graph for "testdata/graph_invalid_missing_deps": dependency error for "one.pgsql": missing file "missing.pgsql"`,
		func() {
			gg.GraphDirInit(`testdata/graph_invalid_missing_deps`)
		},
	)
}

func TestGraphDir_invalid_multiple_entries(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(
		`unable to build dependency graph for "testdata/graph_invalid_multiple_entries": expected to find exactly one dependency-free entry file, found multiple: ["one.pgsql" "two.pgsql"]`,
		func() {
			gg.GraphDirInit(`testdata/graph_invalid_multiple_entries`)
		},
	)
}

func TestGraphDir_invalid_cyclic_self(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(
		`unable to build dependency graph for "testdata/graph_invalid_cyclic_self": dependency cycle: ["one.pgsql" "one.pgsql"]`,
		func() {
			gg.GraphDirInit(`testdata/graph_invalid_cyclic_self`)
		},
	)
}

func TestGraphDir_invalid_cyclic_direct(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(
		`unable to build dependency graph for "testdata/graph_invalid_cyclic_direct": dependency cycle: ["one.pgsql" "two.pgsql" "one.pgsql"]`,
		func() {
			gg.GraphDirInit(`testdata/graph_invalid_cyclic_direct`)
		},
	)
}

func TestGraphDir_invalid_cyclic_indirect(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(
		`unable to build dependency graph for "testdata/graph_invalid_cyclic_indirect": dependency cycle: ["four.pgsql" "one.pgsql" "two.pgsql" "three.pgsql" "four.pgsql"]`,
		func() {
			gg.GraphDirInit(`testdata/graph_invalid_cyclic_indirect`)
		},
	)
}

func TestGraphDir_valid_empty(t *testing.T) {
	defer gtest.Catch(t)

	testGraphDir(`testdata/empty`, nil)
}

func TestGraphDir_valid_non_empty(t *testing.T) {
	defer gtest.Catch(t)

	testGraphDir(`testdata/graph_valid_non_empty`, []string{
		`schema.pgsql`,
		`one.pgsql`,
		`two.pgsql`,
		`three.pgsql`,
		`four.pgsql`,
	})
}

func testGraphDir(dir string, exp []string) {
	gtest.Equal(gg.GraphDirInit(dir).Names(), exp)
}

var graphFileSrc = gg.ReadFile[string](`testdata/graph_file_long`)

var graphFileOut = []string{
	`aaa7c30c9fe6494db244df541a415b8f`,
	`ed33e824fe574a2f91712c1a1609df8c`,
	`ebe9816ee8b14ce9bba478c8e0853581`,
	`b6728a7d157e4984afb430ed2bf750b7`,
	`f4f68f8f00dd45fcba1b2a97c1eafc94`,
	`5acde9df2bb348d1aeb55dbc8f06565c`,
	`e6a34f990e2c4bbd85b13f46d96ed708`,
	`889b367cd42d42189a1b7d9d3f177e84`,
	`00ef58a6eca448c799d744ba5630fc48`,
	`b737450984cd4daea11170364773e98c`,
	`fb37e2f97f3f469080eacd08e29e99ad`,
	`09c3e5a78bf14e69b61b5c8b10db0bec`,
	`e9dd168029cd441296ac6d918c8a95b5`,
	`a83e48bad3eb414c89479bb6666b1e76`,
	`d3316aeb511a4d9295f4b78a3e330bdc`,
	`dac680dcf3fd4f0b99d0789cf396f777`,
	`42d2a4fb764445818d07e5fee726448d`,
}

func Test_graph_file_parse_regexp(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Equal(graphFileParseRegexp(graphFileSrc), graphFileOut)
}

func Benchmark_graph_file_parse_regexp(b *testing.B) {
	defer gtest.Catch(b)

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(graphFileParseRegexp(graphFileSrc))
	}
}

// Copied from `file_graph.go` because it's private.
func graphFileParseRegexp(src string) []string {
	return firstSubmatches(graphFileImportRegexp, src)
}

var graphFileImportRegexp = regexp.MustCompile(`(?m)^@import\s+(.*)$`)

func firstSubmatches(reg *regexp.Regexp, src string) []string {
	return gg.Map(reg.FindAllStringSubmatch(src, -1), get1)
}

func get1(src []string) string { return src[1] }

func Test_graph_file_parse_custom(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Equal(graphFileParseCustom(graphFileSrc), graphFileOut)
}

/*
On the author's machine in Go 1.20.2:

	Benchmark_graph_file_parse_regexp  58900 ns/op  37 allocs/op
	Benchmark_graph_file_parse_custom   8495 ns/op   6 allocs/op

In practice, the difference tends to be significantly lower, making the custom
parser not worth it.
*/
func Benchmark_graph_file_parse_custom(b *testing.B) {
	defer gtest.Catch(b)

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(graphFileParseCustom(graphFileSrc))
	}
}

func graphFileParseCustom(src string) []string {
	var par graphFileParser
	par.src = src
	par.run()
	return par.mat
}

type graphFileParser struct {
	src string
	pos int
	mat []string
}

func (self *graphFileParser) run() {
	if self.more() {
		self.scanImport()
	}

	for self.more() {
		if self.scanned((*graphFileParser).scanNewline) {
			if self.more() {
				self.scanImport()
			}
			continue
		}
		self.skip()
	}
}

func (self *graphFileParser) scanNewline() {
	head := self.headByte()

	if head == '\r' {
		self.pos++
		if self.more() {
			head = self.headByte()
			goto lineFeed
		}
		return
	}

lineFeed:
	if head == '\n' {
		self.pos++
	}
}

func (self *graphFileParser) scanImport() {
	if !self.scannedString(`@import `) {
		return
	}

	var ind int
	var char rune
	rest := self.rest()
	for ind, char = range rest {
		if char == '\r' || char == '\n' {
			break
		}
	}
	self.pos += ind
	self.mat = append(self.mat, strings.TrimSpace(rest[:ind]))
}

func (self *graphFileParser) more() bool { return self.pos < len(self.src) }

func (self *graphFileParser) headByte() byte { return self.src[self.pos] }

func (self *graphFileParser) rest() string {
	if self.more() {
		return self.src[self.pos:]
	}
	return ``
}

func (self *graphFileParser) scanned(fun func(*graphFileParser)) bool {
	sta := self.pos
	if fun != nil {
		fun(self)
	}
	return self.pos > sta
}

func (self *graphFileParser) scannedString(pre string) bool {
	if strings.HasPrefix(self.rest(), pre) {
		self.pos += len(pre)
		return true
	}
	return false
}

/*
Caution: possible antipattern becaues this skips a byte, not a character.
Implemented this way for benchmark reasons. Our public implementation doesn't
do that.
*/
func (self *graphFileParser) skip() { self.pos++ }
