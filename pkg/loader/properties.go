package loader

import (
	"fmt"
	"sort"

	"github.com/stoewer/go-strcase"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

// orderedPropertyKeys returns the keys of m alphabetically ordered
// keys in required will be placed first
func orderedPropertyKeys(required []string, m map[string]apiextensions.JSONSchemaProps, isResource bool) []string {
	sort.Strings(required)

	if isResource {
		mkeys := make(map[string]struct{})
		for k := range m {
			mkeys[k] = struct{}{}
		}
		for _, special := range []string{"metadata", "kind", "apiVersion"} {
			if !isRequired(special, required) {
				if _, ok := mkeys[special]; ok {
					required = append([]string{special}, required...)
				}
			}
		}
	}

	keys := make([]string, len(m)-len(required))
	i := 0
	for k := range m {
		if !isRequired(k, required) {
			keys[i] = k
			i++
		}
	}
	sort.Strings(keys)
	return append(required, keys...)
}

// isRequired returns true if k is in the required array
func isRequired(k string, required []string) bool {
	for _, r := range required {
		if r == k {
			return true
		}
	}
	return false
}

// getTypeNameAndKey returns the display name of a Schema type.
func getTypeNameAndKey(fieldName string, s apiextensions.JSONSchemaProps) (string, *string) {
	// Recurse if type is array
	if s.Type == "array" {
		typ, key := getTypeNameAndKey(fieldName, *s.Items.Schema)
		return fmt.Sprintf("[]%s", typ), key
	}

	// Recurse if type is map
	if s.Type == "object" && s.AdditionalProperties != nil {
		typ, key := getTypeNameAndKey(fieldName, *s.AdditionalProperties.Schema)
		return fmt.Sprintf("map[string]%s", typ), key
	}

	// Handle complex types
	if s.Type == "object" && s.Properties != nil {
		// TODO(roee88): we don't have type information so we need a reasonable workaround here
		key := fmt.Sprintf("%sSpec", strcase.UpperCamelCase(fieldName))
		if fieldName == "spec" {
			key = strcase.UpperCamelCase(fieldName)
		}
		return key, &key
	}

	// Get the value for primitive types
	value := s.Type
	if s.Format != "" && value == "byte" {
		value = "[]byte"
	}

	return value, nil
}
