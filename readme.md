## Overview

Essential utilities missing from the Go standard library.

Docs: https://pkg.go.dev/github.com/mitranim/gg

Some features:

* Designed for Go 1.18+. Takes massive advantage of generics.
* Errors with stack traces.
* Rich tooling for exception-style error handling. A great reduction in verbosity.
* Testing framework: assertions with descriptive errors and traces.
* Goroutine-local storage (GLS) via dynamically scoped variables.
* SQL connector: decode rows into Go structs.
* Functional programming utilities: map, filter, fold, and more.
* Common-sense generic data types: zero-optionals, true optionals, sets, indexed collections, and more.
* Checked math.
* Various shortcuts for reflection.
* Various shortcuts for manipulating slices.
* Various shortcuts for manipulating maps.
* Various shortcuts for manipulating strings.
* CLI flag parsing.
* Carefully designed for compatibility with standard interfaces and interfaces commonly supported by 3rd parties.
* No over-modularization.
* No external dependencies.

Submodules:

* `gtest`: testing and assertion tools.
* `grepr`: tools for printing Go data structures as Go code.
* `gsql`: SQL tools:
  * Support for scanning SQL rows into Go structs.
  * Support for SQL arrays.

Complemented by other essential libraries:
* [mitranim/gt](https://github.com/mitranim/gt): many primitive types with support for _both_ JSON and SQL; also dash-free UUID, date, interval.
* [mitranim/gr](https://github.com/mitranim/gr): library of HTTP request shortcuts with a builder-style API.
* [mitranim/rd](https://github.com/mitranim/rd): decoder of HTTP requests into structs with support for formdata (URL or body), multipart, as well as JSON. Lets your server transparently support multiple formats.
* [mitranim/rout](https://github.com/mitranim/rout): HTTP router with procedural control flow (extremely rare!)
* [mitranim/goh](https://github.com/mitranim/goh): various utility types which implement `http.Handler`; write "functional"-style endpoints without inventing any new interfaces.
* [mitranim/sqlb](https://github.com/mitranim/sqlb): SQL query builder; simple, fast, safe, composable.
* [mitranim/gax](https://github.com/mitranim/gax): write HTML and XML as plain Go code; type-checked and much faster than templating.
* [mitranim/rf](https://github.com/mitranim/rf): tools and shortcuts for reflection; has a "JIT compiler" of efficient walkers for deep Go structures.

## License

https://unlicense.org
