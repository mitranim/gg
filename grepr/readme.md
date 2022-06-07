Missing feature of the standard library: printing arbitrary inputs as Go code,
with proper spacing and support for multi-line output with indentation. The
name "repr" stands for "representation" and alludes to the Python function with
the same name.

API doc: https://pkg.go.dev/github.com/mitranim/gg/grepr

Example:

```go
package mock

import (
  "fmt"

  "github.com/mitranim/gg"
  "github.com/mitranim/gg/grepr"
)

type Outer struct {
  OuterId   int
  OuterName string
  Embed
  Inner *Inner
}

type Embed struct {
  EmbedId   int
  EmbedName string
}

type Inner struct {
  InnerId   *int
  InnerName *string
}

func main() {
  fmt.Println(grepr.String(Outer{
    OuterName: `outer`,
    Embed:     Embed{EmbedId: 20},
    Inner:     &Inner{
      InnerId:   gg.Ptr(30),
      InnerName: gg.Ptr(`inner`),
    },
  }))

  /**
  mock.Outer{
    OuterName: `outer`,
    Embed: mock.Embed{EmbedId: 20},
    Inner: &mock.Inner{
      InnerId: gg.Ptr(30),
      InnerName: gg.Ptr(`inner`),
    },
  }
  */
}
```
