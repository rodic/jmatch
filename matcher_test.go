package jmatch

import (
	"os"
	"reflect"
	"testing"
)

type CollectorMatcher struct {
	matches []string
}

func (fm *CollectorMatcher) Match(path string, token Token) {
	fm.matches = append(fm.matches, path)
}

func TestMatcherValid(t *testing.T) {
	testCases := []struct {
		name     string
		expected map[string]Token
	}{
		{name: "testdata/valid/nested.json",
			expected: map[string]Token{
				".name":                    newToken(String, "Chris", 2, 13),
				".age":                     newToken(Number, "23", 3, 12),
				".address.city":            newToken(String, "New York", 5, 15),
				".address.country":         newToken(String, "America", 6, 18),
				".friends.[0].name":        newToken(String, "Emily", 10, 17),
				".friends.[0].hobbies.[0]": newToken(String, "biking", 11, 22),
				".friends.[0].hobbies.[1]": newToken(String, "music", 11, 32),
				".friends.[0].hobbies.[2]": newToken(String, "gaming", 11, 41),
				".friends.[1].name":        newToken(String, "John", 14, 17),
				".friends.[1].hobbies.[0]": newToken(String, "soccer", 15, 22),
				".friends.[1].hobbies.[1]": newToken(String, "gaming", 15, 32),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			json, err := os.ReadFile(tc.name)
			if err != nil {
				t.Fatal(err)
			}

			collector := collectorMatcher{
				collection: make(map[string]Token),
			}

			Match(string(json), &collector)

			if !reflect.DeepEqual(collector.collection, tc.expected) {
				t.Errorf("Expected '%v', got '%v' instead\n", tc.expected, collector.collection)
			}
		})
	}
}
