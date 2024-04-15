package parser

import (
	"reflect"
	"testing"

	z "github.com/rodic/jmatch/tokenizer"
)

func TestSuccessParse(t *testing.T) {
	testCases := []struct {
		name     string
		tokens   []z.Token
		expected []ParsingResult
	}{
		// single value
		{name: "'1'",
			tokens: []z.Token{
				z.NewStringToken("1", 1, 1),
			},
			expected: []ParsingResult{
				{Path: ".", Token: z.NewStringToken("1", 1, 1)},
			},
		},

		// objects
		{name: "{}",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewRightBraceToken(1, 1),
			},
			expected: []ParsingResult{},
		},
		{name: "{'a': '1'}",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewStringToken("a", 1, 2),
				z.NewColonToken(1, 3),
				z.NewStringToken("1", 1, 4),
				z.NewRightBraceToken(1, 5),
			},
			expected: []ParsingResult{
				{Path: ".a", Token: z.NewStringToken("1", 1, 4)},
			},
		},
		{name: "{'a.b': '1'}",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewStringToken("a.b", 1, 2),
				z.NewColonToken(1, 3),
				z.NewStringToken("1", 1, 4),
				z.NewRightBraceToken(1, 5),
			},
			expected: []ParsingResult{
				{Path: ".\"a.b\"", Token: z.NewStringToken("1", 1, 4)},
			},
		},
		{name: "{'a b': '1'}",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewStringToken("a b", 1, 2),
				z.NewColonToken(1, 3),
				z.NewStringToken("1", 1, 4),
				z.NewRightBraceToken(1, 5),
			},
			expected: []ParsingResult{
				{Path: ".\"a b\"", Token: z.NewStringToken("1", 1, 4)},
			},
		},
		{name: "{'a': '1', 'b': '2', 'c': 3}",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewStringToken("a", 1, 2),
				z.NewColonToken(1, 3),
				z.NewStringToken("1", 1, 4),
				z.NewCommaToken(1, 5),
				z.NewStringToken("b", 1, 6),
				z.NewColonToken(1, 7),
				z.NewStringToken("2", 1, 8),
				z.NewCommaToken(1, 9),
				z.NewStringToken("c", 1, 10),
				z.NewColonToken(1, 111),
				z.NewStringToken("3", 1, 12),
				z.NewRightBraceToken(1, 5),
			},
			expected: []ParsingResult{
				{Path: ".a", Token: z.NewStringToken("1", 1, 4)},
				{Path: ".b", Token: z.NewStringToken("2", 1, 8)},
				{Path: ".c", Token: z.NewStringToken("3", 1, 12)},
			},
		},
		{name: "{'a': { 'b': { 'c': 3}}}",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewStringToken("a", 1, 2),
				z.NewColonToken(1, 3),
				z.NewLeftBraceToken(1, 4),
				z.NewStringToken("b", 1, 5),
				z.NewColonToken(1, 6),
				z.NewLeftBraceToken(1, 7),
				z.NewStringToken("c", 1, 8),
				z.NewColonToken(1, 9),
				z.NewStringToken("3", 1, 10),
				z.NewRightBraceToken(1, 11),
				z.NewRightBraceToken(1, 12),
				z.NewRightBraceToken(1, 13),
			},
			expected: []ParsingResult{
				{Path: ".a.b.c", Token: z.NewStringToken("3", 1, 10)},
			},
		},

		// arrays
		{name: "['1', '2']",
			tokens: []z.Token{
				z.NewLeftBracketToken(1, 1),
				z.NewStringToken("1", 1, 2),
				z.NewCommaToken(1, 3),
				z.NewStringToken("2", 1, 4),
				z.NewRightBracketToken(1, 5),
			},
			expected: []ParsingResult{
				{Path: ".[0]", Token: z.NewStringToken("1", 1, 2)},
				{Path: ".[1]", Token: z.NewStringToken("2", 1, 4)},
			},
		},
		{name: "['1', ['2']]",
			tokens: []z.Token{
				z.NewLeftBracketToken(1, 1),
				z.NewStringToken("1", 1, 2),
				z.NewCommaToken(1, 3),
				z.NewLeftBracketToken(1, 4),
				z.NewStringToken("2", 1, 5),
				z.NewRightBracketToken(1, 6),
				z.NewRightBracketToken(1, 7),
			},
			expected: []ParsingResult{
				{Path: ".[0]", Token: z.NewStringToken("1", 1, 2)},
				{Path: ".[1][0]", Token: z.NewStringToken("2", 1, 5)},
			},
		},
		{name: "[{}]",
			tokens: []z.Token{
				z.NewLeftBracketToken(1, 1),
				z.NewLeftBraceToken(1, 3),
				z.NewRightBraceToken(1, 5),
				z.NewRightBracketToken(1, 1),
			},
			expected: []ParsingResult{},
		},
		{name: "[{'a': 1}]",
			tokens: []z.Token{
				z.NewLeftBracketToken(1, 1),
				z.NewLeftBraceToken(1, 2),
				z.NewStringToken("a", 1, 3),
				z.NewColonToken(1, 4),
				z.NewNumberToken("1", 1, 5),
				z.NewRightBraceToken(1, 6),
				z.NewRightBracketToken(1, 7),
			},
			expected: []ParsingResult{
				{Path: ".[0].a", Token: z.NewNumberToken("1", 1, 5)},
			},
		},
		{name: "[[[[[[1]]]]]]",
			tokens: []z.Token{
				z.NewLeftBracketToken(1, 1),
				z.NewLeftBracketToken(1, 2),
				z.NewLeftBracketToken(1, 3),
				z.NewLeftBracketToken(1, 4),
				z.NewLeftBracketToken(1, 5),
				z.NewLeftBracketToken(1, 6),
				z.NewStringToken("1", 1, 7),
				z.NewRightBracketToken(1, 8),
				z.NewRightBracketToken(1, 9),
				z.NewRightBracketToken(1, 10),
				z.NewRightBracketToken(1, 11),
				z.NewRightBracketToken(1, 12),
				z.NewRightBracketToken(1, 13),
			},
			expected: []ParsingResult{
				{Path: ".[0][0][0][0][0][0]", Token: z.NewStringToken("1", 1, 7)},
			},
		},

		// mixed
		{name: "{'a': ['1', '2']}",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewStringToken("a", 1, 2),
				z.NewColonToken(1, 3),
				z.NewLeftBracketToken(1, 4),
				z.NewStringToken("1", 1, 5),
				z.NewCommaToken(1, 6),
				z.NewStringToken("2", 1, 7),
				z.NewRightBracketToken(1, 1),
				z.NewRightBraceToken(1, 5),
			},
			expected: []ParsingResult{
				{Path: ".a[0]", Token: z.NewStringToken("1", 1, 5)},
				{Path: ".a[1]", Token: z.NewStringToken("2", 1, 7)},
			},
		},
		{name: "{'a': [[[[[[1]]]]]]}",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewStringToken("a", 1, 2),
				z.NewColonToken(1, 3),
				z.NewLeftBracketToken(1, 4),
				z.NewLeftBracketToken(1, 5),
				z.NewLeftBracketToken(1, 6),
				z.NewLeftBracketToken(1, 7),
				z.NewLeftBracketToken(1, 8),
				z.NewLeftBracketToken(1, 9),
				z.NewStringToken("1", 1, 10),
				z.NewRightBracketToken(1, 11),
				z.NewRightBracketToken(1, 12),
				z.NewRightBracketToken(1, 13),
				z.NewRightBracketToken(1, 14),
				z.NewRightBracketToken(1, 15),
				z.NewRightBracketToken(1, 16),
				z.NewRightBraceToken(1, 17),
			},
			expected: []ParsingResult{
				{Path: ".a[0][0][0][0][0][0]", Token: z.NewStringToken("1", 1, 10)},
			},
		},
		{name: "{'a': ['1', ['2', '3']]}",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewStringToken("a", 1, 2),
				z.NewColonToken(1, 3),
				z.NewLeftBracketToken(1, 4),
				z.NewStringToken("1", 1, 5),
				z.NewCommaToken(1, 6),
				z.NewLeftBracketToken(1, 7),
				z.NewStringToken("2", 1, 8),
				z.NewCommaToken(1, 9),
				z.NewStringToken("3", 1, 10),
				z.NewRightBracketToken(1, 11),
				z.NewRightBracketToken(1, 12),
				z.NewRightBraceToken(1, 13),
			},
			expected: []ParsingResult{
				{Path: ".a[0]", Token: z.NewStringToken("1", 1, 5)},
				{Path: ".a[1][0]", Token: z.NewStringToken("2", 1, 8)},
				{Path: ".a[1][1]", Token: z.NewStringToken("3", 1, 10)},
			},
		},
		{name: "{'a': ['1', {'b': ['2', {'c': 3}]}]}",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewStringToken("a", 1, 2),
				z.NewColonToken(1, 3),
				z.NewLeftBracketToken(1, 4),
				z.NewStringToken("1", 1, 5),
				z.NewCommaToken(1, 6),
				z.NewLeftBraceToken(1, 7),
				z.NewStringToken("b", 1, 8),
				z.NewColonToken(1, 9),
				z.NewLeftBracketToken(1, 10),
				z.NewStringToken("2", 1, 11),
				z.NewCommaToken(1, 12),
				z.NewLeftBraceToken(1, 13),
				z.NewStringToken("c", 1, 14),
				z.NewColonToken(1, 15),
				z.NewStringToken("3", 1, 16),
				z.NewRightBraceToken(1, 17),
				z.NewRightBracketToken(1, 18),
				z.NewRightBraceToken(1, 19),
				z.NewRightBracketToken(1, 20),
				z.NewRightBraceToken(1, 21),
			},
			expected: []ParsingResult{
				{Path: ".a[0]", Token: z.NewStringToken("1", 1, 5)},
				{Path: ".a[1].b[0]", Token: z.NewStringToken("2", 1, 11)},
				{Path: ".a[1].b[1].c", Token: z.NewStringToken("3", 1, 16)},
			},
		},
		{name: "{'a': ['1', '2'], 'b': '3', 'c': ['4']}",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewStringToken("a", 1, 2),
				z.NewColonToken(1, 3),
				z.NewLeftBracketToken(1, 4),
				z.NewStringToken("1", 1, 5),
				z.NewCommaToken(1, 6),
				z.NewStringToken("2", 1, 7),
				z.NewRightBracketToken(1, 8),
				z.NewCommaToken(1, 9),
				z.NewStringToken("b", 1, 10),
				z.NewColonToken(1, 11),
				z.NewStringToken("3", 1, 12),
				z.NewCommaToken(1, 13),
				z.NewStringToken("c", 1, 14),
				z.NewColonToken(1, 15),
				z.NewLeftBracketToken(1, 16),
				z.NewStringToken("4", 1, 17),
				z.NewRightBracketToken(1, 18),
				z.NewRightBraceToken(1, 19),
			},
			expected: []ParsingResult{
				{Path: ".a[0]", Token: z.NewStringToken("1", 1, 5)},
				{Path: ".a[1]", Token: z.NewStringToken("2", 1, 7)},
				{Path: ".b", Token: z.NewStringToken("3", 1, 12)},
				{Path: ".c[0]", Token: z.NewStringToken("4", 1, 17)},
			},
		},
		{name: "{'s': {'t': [[1], -2.0, '3', true, {'x': false}]}}",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewStringToken("s", 1, 2),
				z.NewColonToken(1, 3),
				z.NewLeftBraceToken(1, 4),
				z.NewStringToken("t", 1, 5),
				z.NewColonToken(1, 6),
				z.NewLeftBracketToken(1, 7),
				z.NewLeftBracketToken(1, 8),
				z.NewNumberToken("1", 1, 9),
				z.NewRightBracketToken(1, 10),
				z.NewCommaToken(1, 11),
				z.NewNumberToken("-2.0", 1, 12),
				z.NewCommaToken(1, 13),
				z.NewStringToken("3", 1, 14),
				z.NewCommaToken(1, 15),
				z.NewBooleanToken("true", 1, 16),
				z.NewCommaToken(1, 17),
				z.NewLeftBraceToken(1, 18),
				z.NewStringToken("x", 1, 19),
				z.NewColonToken(1, 20),
				z.NewBooleanToken("false", 1, 21),
				z.NewRightBraceToken(1, 22),
				z.NewRightBracketToken(1, 23),
				z.NewRightBraceToken(1, 24),
				z.NewRightBraceToken(1, 25),
			},
			expected: []ParsingResult{
				{Path: ".s.t[0][0]", Token: z.NewNumberToken("1", 1, 9)},
				{Path: ".s.t[1]", Token: z.NewNumberToken("-2.0", 1, 12)},
				{Path: ".s.t[2]", Token: z.NewStringToken("3", 1, 14)},
				{Path: ".s.t[3]", Token: z.NewBooleanToken("true", 1, 16)},
				{Path: ".s.t[4].x", Token: z.NewBooleanToken("false", 1, 21)},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			tokenStream := make(chan z.TokenResult)

			go func() {
				for _, t := range tc.tokens {
					tokenStream <- z.TokenResult{Token: t}
				}
				close(tokenStream)
			}()

			p, err := NewParser(tokenStream)

			if err != nil {
				t.Error(err)
			}

			go p.Parse()

			result := make([]ParsingResult, 0, 10)

			for pr := range p.GetResultReadStream() {
				result = append(result, pr)
			}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected '%v', got '%v' instead\n", tc.expected, result)
			}
		})
	}

}

func TestFailParse(t *testing.T) {
	testCases := []struct {
		name     string
		tokens   []z.Token
		expected string
	}{
		// objects
		{name: "{",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "{{",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewLeftBraceToken(1, 2),
			},
			expected: "invalid JSON. unexpected token { at line 1 column 2"},
		{name: "{{{",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewLeftBraceToken(1, 2),
				z.NewLeftBraceToken(1, 3),
			},
			expected: "invalid JSON. unexpected token { at line 1 column 2"},
		{name: "{{}",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewLeftBraceToken(1, 2),
				z.NewRightBraceToken(1, 3),
			},
			expected: "invalid JSON. unexpected token { at line 1 column 2"},
		{name: "{}}",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewRightBraceToken(1, 2),
				z.NewRightBraceToken(1, 3),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "{1",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewNumberToken("1", 1, 2),
			},
			expected: "invalid JSON. unexpected token 1 at line 1 column 2"},
		{name: "{true",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewBooleanToken("true", 1, 2),
			},
			expected: "invalid JSON. unexpected token true at line 1 column 2"},
		{name: "{null",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewNullToken(1, 2),
			},
			expected: "invalid JSON. unexpected token null at line 1 column 2"},
		{name: "{:",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewColonToken(1, 2),
			},
			expected: "invalid JSON. unexpected token : at line 1 column 2"},
		{name: "{,",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewCommaToken(1, 2),
			},
			expected: "invalid JSON. unexpected token , at line 1 column 2"},
		{name: "{'a',",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewStringToken("a", 1, 2),
				z.NewCommaToken(1, 3),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "{'a',1}",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewStringToken("a", 1, 2),
				z.NewCommaToken(1, 3),
				z.NewNumberToken("1", 1, 4),
				z.NewRightBraceToken(1, 5),
			},
			expected: "invalid JSON. unexpected token , at line 1 column 3"},
		{name: "{'a','1'}",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewStringToken("a", 1, 2),
				z.NewCommaToken(1, 3),
				z.NewStringToken("1", 1, 4),
				z.NewRightBraceToken(1, 5),
			},
			expected: "invalid JSON. unexpected token , at line 1 column 3"},
		{name: "{'a': 1,}",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewStringToken("a", 1, 2),
				z.NewColonToken(1, 3),
				z.NewNumberToken("1", 1, 4),
				z.NewCommaToken(1, 5),
				z.NewRightBraceToken(1, 6),
			},
			expected: "invalid JSON. unexpected token } at line 1 column 6"},
		{name: "{'a': 1, 2}",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewStringToken("a", 1, 2),
				z.NewColonToken(1, 3),
				z.NewNumberToken("1", 1, 4),
				z.NewCommaToken(1, 5),
				z.NewNumberToken("2", 1, 6),
				z.NewRightBraceToken(1, 7),
			},
			expected: "invalid JSON. unexpected token 2 at line 1 column 6"},
		{name: "{'a': 'b': 1}",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewStringToken("a", 1, 2),
				z.NewColonToken(1, 3),
				z.NewStringToken("b", 1, 4),
				z.NewColonToken(1, 5),
				z.NewNumberToken("1", 1, 6),
				z.NewRightBraceToken(1, 7),
			},
			expected: "invalid JSON. unexpected token : at line 1 column 5"},
		{name: "{'a': {,}}",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewStringToken("a", 1, 2),
				z.NewColonToken(1, 3),
				z.NewLeftBraceToken(1, 4),
				z.NewCommaToken(1, 5),
				z.NewRightBraceToken(1, 6),
				z.NewRightBraceToken(1, 7),
			},
			expected: "invalid JSON. unexpected token , at line 1 column 5"},
		{name: "{'a': {{}}}",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewStringToken("a", 1, 2),
				z.NewColonToken(1, 3),
				z.NewLeftBraceToken(1, 4),
				z.NewLeftBraceToken(1, 5),
				z.NewRightBraceToken(1, 6),
				z.NewRightBraceToken(1, 7),
				z.NewRightBraceToken(1, 8),
			},
			expected: "invalid JSON. unexpected token { at line 1 column 5"},

		// arrays
		{name: "[",
			tokens: []z.Token{
				z.NewLeftBracketToken(1, 1),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "[[",
			tokens: []z.Token{
				z.NewLeftBracketToken(1, 1),
				z.NewLeftBracketToken(1, 2),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "[[[",
			tokens: []z.Token{
				z.NewLeftBracketToken(1, 1),
				z.NewLeftBracketToken(1, 2),
				z.NewLeftBracketToken(1, 3),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "[[]",
			tokens: []z.Token{
				z.NewLeftBracketToken(1, 1),
				z.NewLeftBracketToken(1, 2),
				z.NewRightBracketToken(1, 3),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "[]]",
			tokens: []z.Token{
				z.NewLeftBracketToken(1, 1),
				z.NewRightBracketToken(1, 2),
				z.NewRightBracketToken(1, 3),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "[,",
			tokens: []z.Token{
				z.NewLeftBracketToken(1, 1),
				z.NewCommaToken(1, 2),
			},
			expected: "invalid JSON. unexpected token , at line 1 column 2"},
		{name: "[1,]",
			tokens: []z.Token{
				z.NewLeftBracketToken(1, 1),
				z.NewNumberToken("1", 1, 2),
				z.NewCommaToken(1, 3),
				z.NewRightBracketToken(1, 4),
			},
			expected: "invalid JSON. unexpected token ] at line 1 column 4"},

		// mixed
		{name: "{]",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewRightBracketToken(1, 2),
			},
			expected: "invalid JSON. unexpected token ] at line 1 column 2"},
		{name: "[}",
			tokens: []z.Token{
				z.NewLeftBracketToken(1, 1),
				z.NewRightBraceToken(1, 2),
			},
			expected: "invalid JSON. unexpected token } at line 1 column 2"},
		{name: "{[}",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewLeftBracketToken(1, 2),
				z.NewRightBraceToken(1, 3),
			},
			expected: "invalid JSON. unexpected token [ at line 1 column 2"},
		{name: "{'a': [}",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewStringToken("a", 1, 2),
				z.NewColonToken(1, 3),
				z.NewLeftBracketToken(1, 4),
				z.NewRightBraceToken(1, 5),
			},
			expected: "invalid JSON. unexpected token } at line 1 column 5"},
		{name: "{'a': [,",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewStringToken("a", 1, 2),
				z.NewColonToken(1, 3),
				z.NewLeftBracketToken(1, 4),
				z.NewCommaToken(1, 5),
			},
			expected: "invalid JSON. unexpected token , at line 1 column 5"},
		{name: "{'a': [,",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewStringToken("a", 1, 2),
				z.NewColonToken(1, 3),
				z.NewLeftBracketToken(1, 4),
				z.NewCommaToken(1, 5),
			},
			expected: "invalid JSON. unexpected token , at line 1 column 5"},
		{name: "{'a': [{},",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewStringToken("a", 1, 2),
				z.NewColonToken(1, 3),
				z.NewLeftBracketToken(1, 4),
				z.NewLeftBraceToken(1, 5),
				z.NewRightBraceToken(1, 6),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "{'s': {'t': [[1], -2.0, '3', true, {'x': false,}]}}",
			tokens: []z.Token{
				z.NewLeftBraceToken(1, 1),
				z.NewStringToken("s", 1, 2),
				z.NewColonToken(1, 3),
				z.NewLeftBraceToken(1, 4),
				z.NewStringToken("t", 1, 5),
				z.NewColonToken(1, 6),
				z.NewLeftBracketToken(1, 7),
				z.NewLeftBracketToken(1, 8),
				z.NewNumberToken("1", 1, 9),
				z.NewRightBracketToken(1, 10),
				z.NewCommaToken(1, 11),
				z.NewNumberToken("-2.0", 1, 12),
				z.NewCommaToken(1, 13),
				z.NewStringToken("3", 1, 14),
				z.NewCommaToken(1, 15),
				z.NewBooleanToken("true", 1, 16),
				z.NewCommaToken(1, 17),
				z.NewLeftBraceToken(1, 18),
				z.NewStringToken("x", 1, 19),
				z.NewColonToken(1, 20),
				z.NewBooleanToken("false", 1, 21),
				z.NewCommaToken(1, 22),
				z.NewRightBraceToken(1, 23),
				z.NewRightBracketToken(1, 24),
				z.NewRightBraceToken(1, 25),
				z.NewRightBraceToken(1, 26),
			},
			expected: "invalid JSON. unexpected token } at line 1 column 23",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			tokenStream := make(chan z.TokenResult)

			go func() {
				for _, t := range tc.tokens {
					tokenStream <- z.TokenResult{Token: t}
				}
				close(tokenStream)
			}()

			p, err := NewParser(tokenStream)

			if err != nil {
				t.Error(err)
			}

			go p.Parse()

			var lastResult ParsingResult

			for pr := range p.GetResultReadStream() {
				lastResult = pr
			}

			if !reflect.DeepEqual(lastResult.Error.Error(), tc.expected) {
				t.Errorf("Expected '%v', got '%v' instead\n", tc.expected, lastResult.Error)
			}
		})
	}
}
