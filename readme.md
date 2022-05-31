## Overview

Essential utilities missing from the Go standard library.

Docs: https://pkg.go.dev/github.com/mitranim/gg

Some features:

  * Designed for Go 1.18. Takes massive advantage of generics.
  * Errors with stack traces.
  * Various functional programming utilities: map, filter, fold, and more.
  * Various shortcuts for reflection.
  * Various shortcuts for manipulating slices.
  * Various shortcuts for manipulating maps.
  * Various shortcuts for manipulating strings.
  * Common-sense generic data types: zero optionals, true optionals, sets, indexed collections, and more.
  * Various utilities for exception-style error handling, using `panic` and `recover`.
  * CLI flag parsing.
  * Various utilities for testing.
    * Assertion shortcuts with descriptive errors and full stack traces.
  * Carefully designed for compatibility with standard interfaces and interfaces commonly supported by 3rd parties.
  * No over-modularization.
  * No external dependencies.

Submodules:

* `gsql`: SQL tools:
  * Support for scanning SQL rows into Go structs.
  * Support for SQL arrays.
* `gtest`: testing and assertion tools.
* `grepr`: tools for printing Go data structures as Go code.

Current limitations:

  * Not fully documented.

## License

https://unlicense.org
