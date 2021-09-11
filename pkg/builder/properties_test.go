// Copyright 2021 IBM Corp.
// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func TestSortPropertiesByRequired(t *testing.T) {
	shuffle := func(s []string) {
		rand.Shuffle(len(s), func(i, j int) {
			s[i], s[j] = s[j], s[i]
		})
	}

	sorter := propertiesByRequired{}

	// No required; sorts alphabetically.
	sorter.properties = []string{"a", "b", "c", "d"}
	shuffle(sorter.properties)
	sort.Sort(sorter)

	if !reflect.DeepEqual(sorter.properties, []string{"a", "b", "c", "d"}) {
		t.Fatalf("wrong order, got %v", sorter.properties)
	}

	// Required first.
	sorter.required = []string{"b", "d"}
	shuffle(sorter.properties)
	sort.Sort(sorter)

	if !reflect.DeepEqual(sorter.properties, []string{"b", "d", "a", "c"}) {
		t.Fatalf("wrong order, got %v", sorter.properties)
	}

	// All required; sorts alphabetically.
	sorter.required = []string{"a", "b", "c", "d"}
	shuffle(sorter.properties)
	sort.Sort(sorter)

	if !reflect.DeepEqual(sorter.properties, []string{"a", "b", "c", "d"}) {
		t.Fatalf("wrong order, got %v", sorter.properties)
	}
}
