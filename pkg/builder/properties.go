// Copyright 2021 IBM Corp.
// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"sort"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

type propertiesByRequired struct {
	properties []string
	required   []string // must be sorted already
}

func (s propertiesByRequired) Len() int {
	return len(s.properties)
}

func (s propertiesByRequired) Swap(i, j int) {
	s.properties[i], s.properties[j] = s.properties[j], s.properties[i]
}

func (s propertiesByRequired) Less(i, j int) bool {
	left := s.properties[i]
	right := s.properties[j]

	leftRequired := false
	if k := sort.SearchStrings(s.required, left); k < len(s.required) && s.required[k] == left {
		leftRequired = true
	}

	rightRequired := false
	if k := sort.SearchStrings(s.required, right); k < len(s.required) && s.required[k] == right {
		rightRequired = true
	}

	if leftRequired && !rightRequired {
		return true
	}
	if !leftRequired && rightRequired {
		return false
	}
	return left < right
}

// orderedPropertyKeys returns the keys of m alphabetically ordered
// keys in required will be placed first
func orderedPropertyKeys(required []string, m map[string]apiextensions.JSONSchemaProps, isResource bool) []string {
	special := []string{"apiVersion", "kind", "metadata"}

	// All keys excluding special ones in resources
	keys := make([]string, 0, len(m))
	for k := range m {
		if !isResource || !containsString(k, special) {
			keys = append(keys, k)
		}
	}

	// Sort properties alphabetically but place required properties first
	sort.Strings(required) // must be sorted when used in propertiesByRequired
	sortableProperties := propertiesByRequired{properties: keys, required: required}
	sort.Sort(sortableProperties)
	properties := sortableProperties.properties

	return properties
}

func getEnrichedProperty(schema *apiextensions.JSONSchemaProps, fieldName string) apiextensions.JSONSchemaProps {
	property := schema.Properties[fieldName]

	// Special case support for single allOf, anyOf, oneOf
	// TODO: consider adding support for not and for length greater than 1
	var validationProperty *apiextensions.JSONSchemaProps
	if len(schema.AllOf) == 1 {
		validationProperty = getProperty(&schema.AllOf[0], fieldName)
	} else if len(schema.AnyOf) == 1 {
		validationProperty = getProperty(&schema.AnyOf[0], fieldName)
	} else if len(schema.OneOf) == 1 {
		validationProperty = getProperty(&schema.OneOf[0], fieldName)
	}

	if validationProperty != nil {
		// does not set description, type, default, additionalProperties, nullable within an allOf, anyOf, oneOf or not
		// with the exception of the two pattern for x-kubernetes-int-or-string: true (see below).
		if property.Format == "" {
			property.Format = validationProperty.Format
		}
		if property.Title == "" {
			property.Title = validationProperty.Title
		}
		if property.Maximum == nil {
			property.Maximum = validationProperty.Maximum
			property.ExclusiveMaximum = validationProperty.ExclusiveMaximum
		}
		if property.Minimum == nil {
			property.Minimum = validationProperty.Minimum
			property.ExclusiveMinimum = validationProperty.ExclusiveMinimum
		}
		if property.MaxLength == nil {
			property.MaxLength = validationProperty.MaxLength
		}
		if property.MinLength == nil {
			property.MinLength = validationProperty.MinLength
		}
		if property.Pattern == "" {
			property.Pattern = validationProperty.Pattern
		}
		if property.MaxItems == nil {
			property.MaxItems = validationProperty.MaxItems
		}
		if property.MinItems == nil {
			property.MinItems = validationProperty.MinItems
		}
		if property.MultipleOf == nil {
			property.MultipleOf = validationProperty.MultipleOf
		}
		if property.Enum == nil {
			property.Enum = validationProperty.Enum
		}
		if property.MaxProperties == nil {
			property.MaxProperties = validationProperty.MaxProperties
		}
		if property.MinProperties == nil {
			property.MinProperties = validationProperty.MinProperties
		}
		if property.Required == nil {
			property.Required = validationProperty.Required
		}
		if property.XValidations == nil {
			property.XValidations = validationProperty.XValidations
		}

	}

	return property
}

func getProperty(schema *apiextensions.JSONSchemaProps, fieldName string) *apiextensions.JSONSchemaProps {
	if schema != nil {
		property, exists := schema.Properties[fieldName]
		if exists {
			return &property
		}
	}
	return nil
}

func containsString(str string, list []string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}
