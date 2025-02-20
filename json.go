package gg

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Uses `json.Marshal` to encode the given value as JSON, panicking on error.
func JsonEncode[Out Text, Src any](src Src) Out {
	return Try1(JsonEncodeCatch[Out](src))
}

/*
Uses `json.MarshalIndent` to encode the given value as JSON with indentation
controlled by the `Indent` variable, panicking on error.
*/
func JsonEncodeIndent[Out Text, Src any](src Src) Out {
	return Try1(JsonEncodeIndentCatch[Out](src))
}

/*
Same as `json.Marshal` but sometimes marginally more efficient. Avoids spurious
heap escape of the input. May be redundant in later Go versions.
*/
func JsonEncodeCatch[Out Text, Src any](src Src) (Out, error) {
	out, err := json.Marshal(AnyNoEscUnsafe(src))
	return ToText[Out](out), err
}

/*
Same as `json.MarshalIndent`, but uses the default indentation controlled by the
`Indent` variable. Also sometimes marginally more efficient. Avoids spurious
heap escape of the input.
*/
func JsonEncodeIndentCatch[Out Text, Src any](src Src) (Out, error) {
	out, err := json.MarshalIndent(AnyNoEscUnsafe(src), ``, Indent)
	return ToText[Out](out), err
}

// Shortcut for `JsonEncode` for `[]byte`.
func JsonBytes[A any](src A) []byte { return JsonEncode[[]byte](src) }

// Shortcut for `JsonEncodeIndent` for `[]byte`.
func JsonBytesIndent[A any](src A) []byte { return JsonEncodeIndent[[]byte](src) }

// Shortcut for `JsonEncodeCatch` for `[]byte`.
func JsonBytesCatch[A any](src A) ([]byte, error) { return JsonEncodeCatch[[]byte](src) }

// Shortcut for `JsonEncodeIndentCatch` for `[]byte`.
func JsonBytesIndentCatch[A any](src A) ([]byte, error) { return JsonEncodeIndentCatch[[]byte](src) }

/*
Shortcut for implementing JSON encoding of `Nullable` types.
Mostly for internal use.
*/
func JsonBytesNullCatch[A any, B NullableValGetter[A]](val B) ([]byte, error) {
	if val.IsNull() {
		return ToBytes(`null`), nil
	}
	return JsonBytesCatch(val.Get())
}

// Shortcut for `JsonEncode` for `string`.
func JsonString[A any](src A) string { return JsonEncode[string](src) }

// Shortcut for `JsonEncodeIndent` for `string`.
func JsonStringIndent[A any](src A) string { return JsonEncodeIndent[string](src) }

// Shortcut for `JsonEncodeCatch` for `string`.
func JsonStringCatch[A any](src A) (string, error) { return JsonEncodeCatch[string](src) }

// Shortcut for `JsonEncodeIndentCatch` for `string`.
func JsonStringIndentCatch[A any](src A) (string, error) { return JsonEncodeIndentCatch[string](src) }

/*
Shortcut for parsing arbitrary text into the given output, panicking on errors.
If the output pointer is nil, does nothing.
*/
func JsonDecode[Out any, Src Text](src Src, out *Out) { Try(JsonDecodeCatch(src, out)) }

/*
Shortcut for parsing the given text into the given output, ignoring errors.
If the output pointer is nil, does nothing.
*/
func JsonDecodeOpt[Out any, Src Text](src Src, out *Out) { Nop1(JsonDecodeCatch(src, out)) }

/*
Shortcut for parsing the given text into a value of the given type, panicking
on errors.
*/
func JsonDecodeTo[Out any, Src Text](src Src) (out Out) {
	Try(JsonDecodeCatch(src, &out))
	return
}

/*
Shortcut for parsing the given text into the given output, ignoring errors.
If the output pointer is nil, does nothing.
*/
func JsonDecodeOptTo[Out any, Src Text](src Src) (out Out) {
	Nop1(JsonDecodeCatch(src, &out))
	return
}

/*
Parses the given text into the given output. Similar to `json.Unmarshal`, but
avoids the overhead of byte-string conversion and spurious escapes. If the
output pointer is nil, does nothing.
*/
func JsonDecodeCatch[Out any, Src Text](src Src, out *Out) error {
	if out != nil {
		return json.Unmarshal(ToBytes(src), AnyNoEscUnsafe(out))
	}
	return nil
}

/*
Shortcut for decoding the content of the given file into a value of the given
type. Panics on error.
*/
func JsonDecodeFileTo[A any](path string) (out A) {
	JsonDecodeFile(path, &out)
	return
}

/*
Shortcut for decoding the content of the given file into a pointer of the given
type. Panics on error.
*/
func JsonDecodeFile[A any](path string, out *A) {
	if out != nil {
		JsonDecodeClose(Try1(os.Open(path)), NoEscUnsafe(out))
	}
}

/*
Shortcut for writing the JSON encoding of the given value to a file at the given
path. Intermediary directories are created automatically. Any existing file is
truncated.
*/
func JsonEncodeFile[A any](path string, src A) {
	MkdirAll(filepath.Dir(path))

	file := Try1(os.Create(path))
	defer file.Close()

	Try(json.NewEncoder(file).Encode(src))
	Try(file.Close())
}

/*
Uses `json.Decoder` to decode one JSON entry/line from the reader, writing to
the given output. Always closes the reader. Panics on errors.
*/
func JsonDecodeClose[A any](src io.ReadCloser, out *A) {
	defer Close(src)
	if out != nil {
		Try(json.NewDecoder(NoEscUnsafe(src)).Decode(AnyNoEscUnsafe(out)))
	}
}

// True if the input is "null" or blank. Ignores whitespace.
func IsJsonEmpty[A Text](val A) bool {
	src := strings.TrimSpace(ToString(val))
	return src == `` || src == `null`
}
