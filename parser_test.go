package jmatch

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		name     string
		tokens   []Token
		expected ParsingResult
	}{
		{name: "{'a': '1'}",
			tokens: []Token{
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "a"},
				{tokenType: Colon, value: ":"},
				{tokenType: String, value: "1"},
				{tokenType: RightBrace, value: "}"},
			},
			expected: ParsingResult{
				".a": Token{tokenType: String, value: "1"},
			},
		},

		{name: "{'a': '1', 'b': '2', 'c': 3}",
			tokens: []Token{
				{tokenType: LeftBrace, value: "{"},

				{tokenType: String, value: "a"},
				{tokenType: Colon, value: ":"},
				{tokenType: String, value: "1"},
				{tokenType: Comma, value: ","},

				{tokenType: String, value: "b"},
				{tokenType: Colon, value: ":"},
				{tokenType: String, value: "2"},
				{tokenType: Comma, value: ","},

				{tokenType: String, value: "c"},
				{tokenType: Colon, value: ":"},
				{tokenType: String, value: "3"},

				{tokenType: RightBrace, value: "}"},
			},
			expected: ParsingResult{
				".a": Token{tokenType: String, value: "1"},
				".b": Token{tokenType: String, value: "2"},
				".c": Token{tokenType: String, value: "3"},
			},
		},

		{name: "{'a': ['1', '2']}",
			tokens: []Token{
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "a"},
				{tokenType: Colon, value: ":"},
				{tokenType: LeftBracket, value: "["},
				{tokenType: String, value: "1"},
				{tokenType: Comma, value: ","},
				{tokenType: String, value: "2"},
				{tokenType: RightBracket, value: "]"},
				{tokenType: RightBrace, value: "}"},
			},
			expected: ParsingResult{
				".a.[0]": Token{tokenType: String, value: "1"},
				".a.[1]": Token{tokenType: String, value: "2"},
			},
		},

		{name: "{'a': ['1', '2'], 'b': '3', 'c': ['4']}",
			tokens: []Token{
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "a"},
				{tokenType: Colon, value: ":"},
				{tokenType: LeftBracket, value: "["},
				{tokenType: String, value: "1"},
				{tokenType: Comma, value: ","},
				{tokenType: String, value: "2"},
				{tokenType: RightBracket, value: "]"},
				{tokenType: String, value: "b"},
				{tokenType: Colon, value: ":"},
				{tokenType: String, value: "3"},
				{tokenType: Comma, value: ","},
				{tokenType: String, value: "c"},
				{tokenType: Colon, value: ":"},
				{tokenType: LeftBracket, value: "["},
				{tokenType: String, value: "4"},
				{tokenType: RightBracket, value: "]"},
				{tokenType: RightBrace, value: "}"},
			},
			expected: ParsingResult{
				".a.[0]": Token{tokenType: String, value: "1"},
				".a.[1]": Token{tokenType: String, value: "2"},
				".b":     Token{tokenType: String, value: "3"},
				".c.[0]": Token{tokenType: String, value: "4"},
			},
		},

		{name: "{'a': ['1', ['2', '3']]}",
			tokens: []Token{
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "a"},
				{tokenType: Colon, value: ":"},
				{tokenType: LeftBracket, value: "["},
				{tokenType: String, value: "1"},
				{tokenType: Comma, value: ","},
				{tokenType: LeftBracket, value: "]"},
				{tokenType: String, value: "2"},
				{tokenType: Comma, value: ","},
				{tokenType: String, value: "3"},
				{tokenType: RightBracket, value: "]"},
				{tokenType: RightBracket, value: "]"},
				{tokenType: RightBrace, value: "}"},
			},
			expected: ParsingResult{
				".a.[0]":     Token{tokenType: String, value: "1"},
				".a.[1].[0]": Token{tokenType: String, value: "2"},
				".a.[1].[1]": Token{tokenType: String, value: "3"},
			},
		},

		{name: "{'a': {'b': {'c': 3}}}",
			tokens: []Token{
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "a"},
				{tokenType: Colon, value: ":"},
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "b"},
				{tokenType: Colon, value: ":"},
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "c"},
				{tokenType: Colon, value: ":"},
				{tokenType: String, value: "3"},
				{tokenType: RightBrace, value: "}"},
				{tokenType: RightBrace, value: "}"},
				{tokenType: RightBrace, value: "}"},
			},
			expected: ParsingResult{
				".a.b.c": Token{tokenType: String, value: "3"},
			},
		},

		{name: "{'a': {'b': ['1', {'c': 2}]}}",
			tokens: []Token{
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "a"},
				{tokenType: Colon, value: ":"},
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "b"},
				{tokenType: Colon, value: ":"},
				{tokenType: LeftBracket, value: "["},
				{tokenType: String, value: "1"},
				{tokenType: Comma, value: ","},
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "c"},
				{tokenType: Colon, value: ":"},
				{tokenType: Number, value: "2"},
				{tokenType: RightBrace, value: "}"},
				{tokenType: RightBracket, value: "]"},
				{tokenType: RightBrace, value: "}"},
				{tokenType: RightBrace, value: "}"},
			},
			expected: ParsingResult{
				".a.b.[0]":   Token{tokenType: String, value: "1"},
				".a.b.[1].c": Token{tokenType: Number, value: "2"},
			},
		},

		{name: "{'s': {'t': [[1], -2.0, '3', true, {'x': false}]}}",
			tokens: []Token{
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "s"},
				{tokenType: Colon, value: ":"},
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "t"},
				{tokenType: Colon, value: ":"},
				{tokenType: LeftBracket, value: "["},
				{tokenType: LeftBracket, value: "["},
				{tokenType: Number, value: "1"},
				{tokenType: RightBracket, value: "]"},
				{tokenType: Comma, value: ","},
				{tokenType: Number, value: "-2.0"},
				{tokenType: Comma, value: ","},
				{tokenType: String, value: "3"},
				{tokenType: Comma, value: ","},
				{tokenType: Boolean, value: "true"},
				{tokenType: Comma, value: ","},
				{tokenType: LeftBrace, value: "{"},
				{tokenType: String, value: "x"},
				{tokenType: Colon, value: ":"},
				{tokenType: Boolean, value: "false"},
				{tokenType: RightBrace, value: "}"},
				{tokenType: RightBracket, value: "]"},
				{tokenType: RightBrace, value: "}"},
				{tokenType: RightBrace, value: "}"},
			},
			expected: ParsingResult{
				".s.t.[0].[0]": Token{tokenType: Number, value: "1"},
				".s.t.[1]":     Token{tokenType: Number, value: "-2.0"},
				".s.t.[2]":     Token{tokenType: String, value: "3"},
				".s.t.[3]":     Token{tokenType: Boolean, value: "true"},
				".s.t.[4].x":   Token{tokenType: Boolean, value: "false"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.tokens)
			paths, err := p.parse()
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(paths, tc.expected) {
				t.Errorf("Expected '%v', got '%v' instead\n", tc.expected, paths)
			}
		})
	}

}
