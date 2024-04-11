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
		{name: "testdata/valid/arrays.json",
			expected: map[string]Token{
				".[0].name": newToken(String, "Chris", 2, 15),
				".[0].age":  newToken(Number, "23", 2, 31),
				".[0].city": newToken(String, "New York", 2, 43),
				".[1].name": newToken(String, "Emily", 3, 15),
				".[1].age":  newToken(Number, "19", 3, 31),
				".[1].city": newToken(String, "Atlanta", 3, 43),
				".[2].name": newToken(String, "Joe", 4, 15),
				".[2].age":  newToken(Number, "32", 4, 29),
				".[2].city": newToken(String, "New York", 4, 41),
				".[3].name": newToken(String, "Kevin", 5, 15),
				".[3].age":  newToken(Number, "19", 5, 31),
				".[3].city": newToken(String, "Atlanta", 5, 43),
				".[4].name": newToken(String, "Michelle", 6, 15),
				".[4].age":  newToken(Number, "27", 6, 34),
				".[4].city": newToken(String, "Los Angeles", 6, 46),
				".[5].name": newToken(String, "Robert", 7, 15),
				".[5].age":  newToken(Number, "45", 7, 32),
				".[5].city": newToken(String, "Manhattan", 7, 44),
				".[6].name": newToken(String, "Sarah", 8, 15),
				".[6].age":  newToken(Number, "31", 8, 31),
				".[6].city": newToken(String, "New York", 8, 43),
			},
		},
		{name: "testdata/valid/date.json",
			expected: map[string]Token{
				".id":      newToken(String, "a98d1377-2270-45fd-8e25-cde720c50bce", 2, 11),
				".message": newToken(String, "Hi Jane 😃 are you busy tonight? Shall we go out for dinner?", 3, 16),
				".date":    newToken(String, "2023-07-24T12:56:15.609Z", 4, 13),
			},
		},
		{name: "testdata/valid/colors.json",
			expected: map[string]Token{
				".[0].calendarId": newToken(String, "e2a5c", 3, 21),
				".[0].color":      newToken(String, "#3997f5", 4, 16),
				".[1].calendarId": newToken(String, "aa027", 7, 21),
				".[1].color":      newToken(String, "#ef5353", 8, 16),
				".[2].calendarId": newToken(String, "5d9a1", 11, 21),
				".[2].color":      newToken(String, "#3fc13f", 12, 16),
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
