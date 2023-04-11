package gg

import (
	"crypto/rand"
	"io"
	"math/big"
)

/*
Generates a random integer of the given type, using the given entropy source
(typically `"crypto/rand".Reader`).
*/
func RandomInt[A Int](src io.Reader) (out A) {
	Try1(src.Read(AsBytes(&out)))
	return out
}

/*
Generates a random integer in the range `[min,max)`, using the given entropy
source (typically `"crypto/rand".Reader`). All numbers must be below
`math.MaxInt64`.
*/
func RandomIntBetween[A Int](src io.Reader, min, max A) A {
	if !(max > min) {
		panic(Errf(`invalid range [%v,%v)`, min, max))
	}

	// The following is suboptimal. See implementation notes below.
	minInt := NumConv[int64](min)
	maxInt := NumConv[int64](max)
	maxBig := new(big.Int).SetInt64(maxInt - minInt)
	tarBig := Try1(rand.Int(src, maxBig))
	return min + NumConv[A](tarBig.Int64())
}

/*
The following implementation doesn't fully pass our test. It performs marginally
better than the wrapper around the "crypto/rand" version. TODO fix. Also TODO
generalize for all int types.

	func RandomUint64Between(src io.Reader, min, max uint64) (out uint64) {
		if !(max > min) {
			panic(Errf(`invalid range [%v,%v)`, min, max))
		}

		ceil := max - min
		bits := bits.Len64(ceil)
		buf := AsBytes(&out)[:(bits+7)/8]

		for {
			Try1(src.Read(buf))
			buf[0] >>= (8 - (bits % 8))

			out = 0
			for _, byte := range buf {
				out = (out << 8) | uint64(byte)
			}
			if out < ceil {
				return out + min
			}
		}
	}
*/

/*
Picks a random element from the given slice, using the given entropy source
(typically `"crypto/rand".Reader`). Panics if the slice is empty or the reader
is nil.
*/
func RandomElem[A any](reader io.Reader, slice []A) A {
	return slice[RandomIntBetween(reader, 0, len(slice))]
}
