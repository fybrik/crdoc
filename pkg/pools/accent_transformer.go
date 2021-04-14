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

package pools

// The code in this file is modified from:
// https://github.com/gohugoio/hugo/blob/d90e37e0c6e812f9913bf256c9c81aa05b7a08aa/common/text/transform.go

import (
	"sync"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var accentTransformerPool = &sync.Pool{
	New: func() interface{} {
		return transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	},
}

// GetAccentsTransformer returns a transformer from the pool.
// The transformer converts text to ascii only.
func GetAccentsTransformer() transform.Transformer {
	return accentTransformerPool.Get().(transform.Transformer)
}

// PutAccentsTransformer returns a transformer to the pool.
// The transformer is reset before it is put back into circulation.
func PutAccentsTransformer(t transform.Transformer) {
	t.Reset()
	accentTransformerPool.Put(t)
}
