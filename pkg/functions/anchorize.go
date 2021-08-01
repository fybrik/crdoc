// Copyright 2015 The Hugo Authors. All rights reserved.
// Modifications copyright (C) 2021 IBM Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package functions

// The code in this file is modified from:
// https://github.com/gohugoio/hugo/blob/d90e37e0c6e812f9913bf256c9c81aa05b7a08aa/markup/goldmark/autoid.go#L43

import (
	"bytes"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/transform"

	"fybrik.io/crdoc/pkg/pools"
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
