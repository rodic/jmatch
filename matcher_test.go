package jmatch

import (
	"reflect"
	"testing"
)

type FixedTokenValueMatch struct {
	matchingString string
	matches        []string
}

func (fm *FixedTokenValueMatch) match(path string, token Token) {
	if token.Value == fm.matchingString {
		fm.matches = append(fm.matches, path)
	}
}

func TestMatcher(t *testing.T) {
	testCases := []struct {
		name     string
		json     string
		expected []string
	}{
		{name: "empty",
			json:     "{}",
			expected: []string{}},
		{name: "simple match",
			json:     "{\"a\": 1}",
			expected: []string{".a"}},
		{name: "multi",
			json:     "{\"a\": 1, \"b\": 2, \"c\": 1}",
			expected: []string{".a", ".c"}},
		{name: "nested",
			json:     "{\"a\": {\"b\": 1}}",
			expected: []string{".a.b"}},
		{name: "nested with arrays",
			json:     "{\"a\": {\"b\": [0, 1, 2]}}",
			expected: []string{".a.b.[1]"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fm := FixedTokenValueMatch{
				matchingString: "1",
				matches:        make([]string, 0, 8),
			}
			_, err := Match(tc.json, &fm)

			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(fm.matches, tc.expected) {
				t.Errorf("Expected '%v', got '%v' instead\n", tc.expected, fm.matches)
			}
		})
	}
}
