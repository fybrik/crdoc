package functions

// The code in this file is modified from:
// https://github.com/gohugoio/hugo/blob/d90e37e0c6e812f9913bf256c9c81aa05b7a08aa/markup/goldmark/autoid.go#L43

import (
	"bytes"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/transform"

	"github.com/mesh-for-data/crdoc/pkg/pools"
)

// Anchorize sanitizes anchor names using GitHub style anchors.
func Anchorize(text string) string {
	return string(anchorize([]byte(text), false))
}

// AnchorizeAsciiOnly sanitizes anchor names using GitHub style anchors
// with non ascii characters stipped out.
func AnchorizeAsciiOnly(text string) string {
	return string(anchorize([]byte(text), true))
}

func anchorize(b []byte, asciiOnly bool) []byte {
	buf := pools.GetBuffer()
	defer pools.PutBuffer(buf)

	if asciiOnly {
		t := pools.GetAccentsTransformer()
		defer pools.PutAccentsTransformer(t)
		b, _, _ = transform.Bytes(t, b)
	}

	b = bytes.TrimSpace(b)

	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		switch {
		case asciiOnly && size != 1:
		case r == '-' || r == ' ':
			buf.WriteRune('-')
		case isAlphaNumeric(r):
			buf.WriteRune(unicode.ToLower(r))
		default:
		}

		b = b[size:]
	}

	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())

	return result
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
