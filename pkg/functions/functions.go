// Copyright 2021 IBM Corp.
// SPDX-License-Identifier: Apache-2.0

package functions

import "html"

var ExportedMap = map[string]interface{}{
	// Strings
	"anchorize":          Anchorize,
	"anchorizeAsciiOnly": AnchorizeAsciiOnly,
	"escape":             html.EscapeString,
}
