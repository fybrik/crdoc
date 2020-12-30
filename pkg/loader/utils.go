package loader

import (
	"regexp"
	"strings"

	"github.com/stoewer/go-strcase"
)

// escapeName returns a name usable as file name
func escapeName(parts ...string) string {
	result := []string{}
	for _, s := range parts {
		if s != "" {
			result = append(result, strcase.KebabCase(s))
		}
	}
	return strings.Join(result, "-")
}

// headingID returns the ID built by hugo for a given header
func headingID(s string) string {
	result := s
	result = strings.ToLower(s)
	result = strings.TrimSpace(result)
	result = regexp.MustCompile(`([^\w\- ]+)`).ReplaceAllString(result, "")
	result = regexp.MustCompile(`(\s)`).ReplaceAllString(result, "-")
	result = regexp.MustCompile(`(\-+$)`).ReplaceAllString(result, "")

	return result
}
