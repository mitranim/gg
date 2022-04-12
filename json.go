package gg

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

// Uses `json.Marshal` to encode the given value as JSON, panicking on error.
func JsonBytes[A any](val A) []byte {
	return Try1(JsonBytesCatch(val))
}

/*
Uses `json.MarshalIndent` to encode the given value as JSON with indentation
controlled by the `Indent` variable, panicking on error.
*/
func JsonBytesIndent[A any](val A) []byte {
	return Try1(JsonBytesIndentCatch(val))
}

/*
Same as `json.Marshal` but sometimes marginally more efficient. Avoids spurious
heap escape of the input.
*/
func JsonBytesCatch[A any](val A) ([]byte, error) {
	return json.Marshal(AnyNoEscUnsafe(val))
}

/*
Same as `json.MarshalIndent`, but uses the default indentation controlled by the
`Indent` variable. Also sometimes marginally more efficient. Avoids spurious
heap escape of the input.
*/
func JsonBytesIndentCatch[A any](val A) ([]byte, error) {
	return json.MarshalIndent(AnyNoEscUnsafe(val), ``, Indent)
}

// Encodes the input as a JSON string, panicking on error.
func JsonString[A any](val A) string { return ToString(JsonBytes(val)) }

/*
Encodes the input as a JSON string, using default indentation controlled by the
`Indent` variable.
*/
func JsonStringIndent[A any](val A) string { return ToString(JsonBytesIndent(val)) }

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

/*
Shortcut for parsing the given string or byte slice into a value of the given
type. Panics on errors.
*/
func JsonParseTo[Out any, Src Text](src Src) (out Out) {
	JsonParse(src, &out)
	return
}

/*
Shortcut for parsing the given string or byte slice into a pointer of the given
type. Panics on errors.
*/
func JsonParse[Out any, Src Text](src Src, out *Out) {
	Try(JsonParseCatch(src, out))
}

/*
Parses the given string or byte slice into a pointer of the given type. Similar
to `json.Unmarshal`, but avoids the overhead of byte-string conversion and
spurious escapes.
*/
func JsonParseCatch[Out any, Src Text](src Src, out *Out) error {
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
