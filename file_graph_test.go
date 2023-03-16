package gg_test

import (
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func TestGraphDir_invalid_missing_deps(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(
		`unable to build dependency graph for "testdata/graph_invalid_missing_deps": dependency error for "one.pgsql": missing file "missing.pgsql"`,
		func() {
			graphDirInit(`testdata/graph_invalid_missing_deps`)
		},
	)
}

func TestGraphDir_invalid_multiple_entries(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(
		`unable to build dependency graph for "testdata/graph_invalid_multiple_entries": expected to find exactly one dependency-free entry file, found multiple: ["one.pgsql" "two.pgsql"]`,
		func() {
			graphDirInit(`testdata/graph_invalid_multiple_entries`)
		},
	)
}

func TestGraphDir_invalid_cyclic_self(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(
		`unable to build dependency graph for "testdata/graph_invalid_cyclic_self": dependency cycle: ["one.pgsql" "one.pgsql"]`,
		func() {
			graphDirInit(`testdata/graph_invalid_cyclic_self`)
		},
	)
}

func TestGraphDir_invalid_cyclic_direct(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(
		`unable to build dependency graph for "testdata/graph_invalid_cyclic_direct": dependency cycle: ["one.pgsql" "two.pgsql" "one.pgsql"]`,
		func() {
			graphDirInit(`testdata/graph_invalid_cyclic_direct`)
		},
	)
}

func TestGraphDir_invalid_cyclic_indirect(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(
		`unable to build dependency graph for "testdata/graph_invalid_cyclic_indirect": dependency cycle: ["four.pgsql" "one.pgsql" "two.pgsql" "three.pgsql" "four.pgsql"]`,
		func() {
			graphDirInit(`testdata/graph_invalid_cyclic_indirect`)
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
		`three.pgsql`,
		`two.pgsql`,
		`four.pgsql`,
	})
}

func testGraphDir(dir string, exp []string) {
	gtest.Equal(graphDirInit(dir).Names(), exp)
}

func graphDirInit(dir string) gg.GraphDir {
	var out gg.GraphDir
	out.Path = dir
	out.Init()
	return out
}
