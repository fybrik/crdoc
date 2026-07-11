// Copyright 2021 IBM Corp.
// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestBuiltinTemplatesRenderStringConstraints is a regression test for the gap
// where the built-in templates rendered minimum/maximum but dropped the sibling
// string constraints minLength, maxLength, and pattern. The fixture puts each
// constraint on its own field so a guard branch cannot be covered by an unrelated
// sibling, and its pattern carries a named capture group ('<', '>'), an entity-like
// sequence ('&'), and an alternation ('|') so the escaping is exercised too:
//   - the HTML templates (markdown, frontmatter) must HTML-escape '<', '>', '&';
//   - the AsciiDoc template must escape '|' as '\|' so it cannot start a new cell.
func TestBuiltinTemplatesRenderStringConstraints(t *testing.T) {
	crds, err := LoadCRDs(filepath.Join("testdata", "constraints.yaml"))
	if err != nil {
		t.Fatalf("LoadCRDs: %v", err)
	}

	cases := []struct {
		template string
		want     []string
		notWant  []string
	}{
		{
			template: "markdown.tmpl",
			want: []string{
				"Min Length</i>: 1",
				"Max Length</i>: 64",
				"Pattern</i>: ^(?P&lt;side&gt;East|West)&amp;copy;$",
			},
			// A raw named-capture group would break HTML rendering; it must be escaped.
			notWant: []string{"?P<side>"},
		},
		{
			template: "frontmatter.tmpl",
			want: []string{
				"Min Length</i>: 1",
				"Max Length</i>: 64",
				"Pattern</i>: ^(?P&lt;side&gt;East|West)&amp;copy;$",
			},
			notWant: []string{"?P<side>"},
		},
		{
			template: "adoc.tmpl",
			want: []string{
				"Min Length_: 1",
				"Max Length_: 64",
				`Pattern_: ^(?P<side>East\|West)&copy;$`,
			},
			// An unescaped '|' would split the pipe-delimited table row.
			notWant: []string{`Pattern_: ^(?P<side>East|West)`},
		},
	}

	for _, tc := range cases {
		t.Run(tc.template, func(t *testing.T) {
			model, err := LoadModel("")
			if err != nil {
				t.Fatalf("LoadModel: %v", err)
			}
			outPath := filepath.Join(t.TempDir(), "out")
			templatePath := filepath.Join("..", "..", "templates", tc.template)
			b := NewModelBuilder(model, false, templatePath, outPath, embed.FS{})
			for _, crd := range crds {
				if err := b.Add(crd); err != nil {
					t.Fatalf("Add: %v", err)
				}
			}
			if err := b.Output(); err != nil {
				t.Fatalf("Output: %v", err)
			}
			data, err := os.ReadFile(outPath)
			if err != nil {
				t.Fatalf("ReadFile: %v", err)
			}
			got := string(data)
			for _, w := range tc.want {
				if !strings.Contains(got, w) {
					t.Errorf("%s: output does not contain %q\n---\n%s", tc.template, w, got)
				}
			}
			for _, nw := range tc.notWant {
				if strings.Contains(got, nw) {
					t.Errorf("%s: output unexpectedly contains %q (escaping regression)\n---\n%s", tc.template, nw, got)
				}
			}
		})
	}
}
