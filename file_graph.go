package gg

import (
	"path/filepath"
	"regexp"
)

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
the dependency graph. After calling this method, the order of files in
`.Files.Slice` represents valid execution order.
*/
func (self *GraphDir) Init() {
	defer Detailf(`unable to build dependency graph for %q`, self.Path)
	self.read()
	self.validateExisting()
	self.validateAcyclic()
	self.sort()
	self.validateEntryFile()
}

// Returns the names of `.Files`, in the same order.
func (self GraphDir) Names() []string {
	return Map(self.Files.Slice, GraphFile.Pk)
}

/*
Returns the `GraphFile` indexed by the given key. Panics if the file is not
found by this key.
*/
func (self GraphDir) File(key string) GraphFile {
	val, ok := self.Files.Got(key)
	if !ok {
		panic(Errf(`missing file %q`, key))
	}
	// Precaution due to use of sorting.
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

func (self *GraphDir) validateExisting() {
	Each(self.Files.Slice, self.validateExistingDeps)
}

func (self *GraphDir) validateExistingDeps(file GraphFile) {
	defer Detailf(`dependency error for %q`, file.Name)

	for _, dep := range file.Deps.Slice {
		Nop1(self.File(dep))
	}
}

func (self *GraphDir) validateAcyclic() {
	var tar graphCycleValidator
	tar.Dir = self
	tar.Run()
}

func (self *GraphDir) sort() {
	Sort(self.Files.Slice)
	self.Files.Reindex()
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
	deps := head.Deps.Len()

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
	Name string
	Body string
	Deps OrdSet[string]
}

// Implement `Pked` for compatibility with `Coll`. See `GraphDir.Files`.
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

// Provisional. Suboptimal.
func (self *GraphFile) parse() {
	self.Deps = OrdSetOf(firstSubmatches(reGraphImport.Get(), self.Body)...)
	self.validateDeps()
}

func (self GraphFile) validateDeps() {
	invalid := Reject(self.Deps.Slice, isBaseName)
	if IsEmpty(invalid) {
		return
	}
	panic(Errf(`invalid imports in %q, every import must be a base name, found %q`, self.Name, invalid))
}

/*
Implement `Lesser` for sorting. Rules:

	* Direct mutual dependency between two files is forbidden. This doesn't forbid
	  indirect cycles, which we validate separately.
	* If one of the files depends on the other, the dependency is "less" and the
	  dependent is "more".
	* If one of the files has fewer dependencies than the other, it's "less" than
	  the other. This also ensures that "pure" dependencies (files which don't
	  depend on others) are "less".
	* Otherwise, determine "less" by comparing file names.
	* If every other condition fails, the receiver is "less" and the input
	  is "more".
*/
func (self GraphFile) Less(other GraphFile) bool {
	has0 := self.Deps.Has(other.Name)
	has1 := other.Deps.Has(self.Name)
	if has0 && has1 {
		panic(Errf(`dependency cycle between %q and %q`))
	}
	if has0 {
		return false
	}
	if has1 {
		return true
	}

	len0 := self.Deps.Len()
	len1 := other.Deps.Len()
	if len0 < len1 {
		return true
	}
	if len1 < len0 {
		return false
	}

	return other.Name > self.Name
}

func (self GraphFile) isEntry() bool { return self.Deps.IsEmpty() }

var reGraphImport = NewLazy(func() *regexp.Regexp {
	return regexp.MustCompile(`(?m)^@import\s+(.*)$`)
})

func isBaseName(val string) bool { return filepath.Base(val) == val }

/*
Forbids cycles. In other words, ensures that our graph is a "multitree".
See https://en.wikipedia.org/wiki/Multitree.
*/
type graphCycleValidator struct {
	Dir   *GraphDir
	valid Set[string]
}

func (self *graphCycleValidator) Run() {
	for _, file := range self.Dir.Files.Slice {
		self.walk(nil, file)
	}
}

func (self *graphCycleValidator) walk(tail *node[string], file GraphFile) {
	if self.valid.Has(file.Name) {
		return
	}

	pending := tail != nil && tail.has(file.Name)
	head := tail.cons(file.Name)

	if pending {
		panic(Errf(`dependency cycle: %q`, Reversed(head.vals())))
	}

	for _, dep := range file.Deps.Slice {
		self.walk(&head, self.Dir.File(dep))
	}
	self.valid.Init().Add(file.Name)
}
