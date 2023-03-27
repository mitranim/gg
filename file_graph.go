package gg

import (
	"path/filepath"
	"regexp"
)

/*
Shortcut for making `GraphDir` with the given path and fully initializing it via
`.Init`.
*/
func GraphDirInit(path string) GraphDir {
	var out GraphDir
	out.Path = path
	out.Init()
	return out
}

/*
Represents a directory where the files form a graph by "importing" each other,
by using special annotations understood by this tool. Supports reading files
from the filesystem, validating the dependency graph, and calculating valid
execution order for the resulting graph. Mostly designed and suited for
emulating a module system for SQL files. May be useful in other similar cases.

The import annotation is currently not customizable and must look like the
following example. Each entry must be placed at the beginning of a line. In
files that contain code, do this within multi-line comments without any prefix.

	@import some_file_name_0
	@import some_file_name_1

Current limitations:

	* The import annotation is non-customizable.
	* No support for file filtering.
	* No support for relative paths. Imports must refer to files by base names.
	* No support for `fs.FS` or other ways to customize reading.
	  Always uses the OS filesystem.
*/
type GraphDir struct {
	Path  string
	Files Coll[string, GraphFile]
}

/*
Reads the files in the directory specified by `.Path`, then builds and validates
the dependency graph. After calling this method, the files in `.Files.Slice`
represent valid execution order.
*/
func (self *GraphDir) Init() {
	defer Detailf(`unable to build dependency graph for %q`, self.Path)
	self.read()
	self.validateExisting()
	self.walk()
	self.validateEntryFile()
}

// Returns the names of `.Files`, in the same order.
func (self GraphDir) Names() []string {
	return Map(self.Files.Slice, GraphFile.Pk)
}

/*
Returns the `GraphFile` indexed by the given key.
Panics if the file is not found.
*/
func (self GraphDir) File(key string) GraphFile {
	val, ok := self.Files.Got(key)
	if !ok {
		panic(Errf(`missing file %q`, key))
	}

	if val.Name != key {
		panic(Errf(`invalid index for %q, found %q`, key, val.Name))
	}

	return val
}

func (self *GraphDir) read() {
	for _, src := range ReadDir(self.Path) {
		if src == nil || src.IsDir() {
			continue
		}

		var file GraphFile
		file.Name = src.Name()
		file.Init(self.Path)

		self.Files.Add(file)
	}
}

// Technically redundant because `graphWalk` also validates this.
func (self GraphDir) validateExisting() {
	Each(self.Files.Slice, self.validateExistingDeps)
}

func (self GraphDir) validateExistingDeps(file GraphFile) {
	defer Detailf(`dependency error for %q`, file.Name)

	for _, dep := range file.Deps {
		Nop1(self.File(dep))
	}
}

func (self *GraphDir) walk() {
	// Forbids cycles and finds valid execution order.
	var walk graphWalk
	walk.Dir = self
	walk.Run()

	// Internal sanity check. If walk is successful, it must build an equivalent
	// set of files. We could also compare the actual elements, but this should
	// be enough to detect mismatches.
	valid := walk.Valid
	len0 := self.Files.Len()
	len1 := valid.Len()
	if len0 != len1 {
		panic(Errf(`internal error: mismatch between original files (length %v) and walked files (length %v)`, len0, len1))
	}

	self.Files = valid
}

/*
Ensures that the resulting graph is either empty, or contains exactly one "entry
file", a file with no dependencies, and that this file has been sorted to the
beginning of the collection. Every other file must explicitly specify its
dependencies. This helps ensure canonical order.
*/
func (self GraphDir) validateEntryFile() {
	if self.Files.IsEmpty() {
		return
	}

	head := Head(self.Files.Slice)
	deps := len(head.Deps)
	if deps != 0 {
		panic(Errf(`expected to begin with a dependency-free entry file, found %q with %v dependencies`, head.Name, deps))
	}

	if None(Tail(self.Files.Slice), GraphFile.isEntry) {
		return
	}

	panic(Errf(
		`expected to find exactly one dependency-free entry file, found multiple: %q`,
		Map(Filter(self.Files.Slice, GraphFile.isEntry), GraphFile.Pk),
	))
}

/*
Represents a file in a graph of files that import each other by using special
import annotations understood by this tool. See `GraphDir` for explanation.
*/
type GraphFile struct {
	Name string   // Valid file base name.
	Body string   // Read from disk by `.Init`.
	Deps []string // Parsed from `.Body` by `.Init`.
}

// Implement `Pker` for compatibility with `Coll`. See `GraphDir.Files`.
func (self GraphFile) Pk() string { return self.Name }

/*
Reads the file named by `.Name` in the given directory, and parses the import
annotations into `.Deps`. Used automatically by `GraphDir.Init`.
*/
func (self *GraphFile) Init(dir string) {
	self.read(dir)
	self.parse()
}

func (self *GraphFile) read(dir string) {
	self.validateName()
	self.Body = ReadFile[string](filepath.Join(dir, self.Name))
}

func (self GraphFile) validateName() {
	if !isBaseName(self.Name) {
		panic(Errf(`unexpected non-base file name %q; file graph currently supports only base-name imports`, self.Name))
	}
}

func (self *GraphFile) parse() {
	/**
	Suboptimal. This is also the slowest part of the `GraphDir` API. Total
	execution can take several milliseconds in some real-life projects with
	hundreds of files. A decently-written custom parser may be able to perform
	several times better.
	*/
	deps := firstSubmatches(reGraphImport.Get(), self.Body)

	invalid := Reject(deps, isBaseName)
	if IsNotEmpty(invalid) {
		panic(Errf(`invalid imports in %q, every import must be a base name, found %q`, self.Name, invalid))
	}

	self.Deps = deps
}

func (self GraphFile) isEntry() bool { return IsEmpty(self.Deps) }

var reGraphImport = NewLazy(func() *regexp.Regexp {
	return regexp.MustCompile(`(?m)^@import\s+(.*)$`)
})

func isBaseName(val string) bool { return filepath.Base(val) == val }

/*
Features:

	* Forbids cycles. In other words, ensures that our graph is a "multitree".
	  See https://en.wikipedia.org/wiki/Multitree.
	* Determines valid execution order.
*/
type graphWalk struct {
	Dir   *GraphDir
	Valid Coll[string, GraphFile]
}

func (self *graphWalk) Run() {
	for _, val := range self.Dir.Files.Slice {
		self.walk(nil, val)
	}
}

func (self *graphWalk) walk(tail *node[string], file GraphFile) {
	key := file.Name
	if self.Valid.Has(key) {
		return
	}

	pending := tail != nil && tail.has(key)
	head := tail.cons(key)

	if pending {
		panic(Errf(`dependency cycle: %q`, Reversed(head.vals())))
	}

	for _, dep := range file.Deps {
		self.walk(&head, self.Dir.File(dep))
	}
	self.Valid.Add(file)
}
