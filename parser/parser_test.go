package parser

import (
	"reflect"
	"testing"

	token "github.com/rodic/jmatch/token"
	z "github.com/rodic/jmatch/tokenizer"
)

func TestSuccessParse(t *testing.T) {
	testCases := []struct {
		name     string
		tokens   []token.Token
		expected []ParsingResult
	}{
		// single value
		{name: "'1'",
			tokens: []token.Token{
				token.NewStringToken("1", 1, 1),
			},
			expected: []ParsingResult{
				{Path: ".", Token: token.NewStringToken("1", 1, 1)},
			},
		},

		// objects
		{name: "{}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewRightBraceToken(1, 1),
			},
			expected: []ParsingResult{},
		},
		{name: "{'a': '1'}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewStringToken("1", 1, 4),
				token.NewRightBraceToken(1, 5),
			},
			expected: []ParsingResult{
				{Path: ".a", Token: token.NewStringToken("1", 1, 4)},
			},
		},
		{name: "{'a.b': '1'}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a.b", 1, 2),
				token.NewColonToken(1, 3),
				token.NewStringToken("1", 1, 4),
				token.NewRightBraceToken(1, 5),
			},
			expected: []ParsingResult{
				{Path: ".\"a.b\"", Token: token.NewStringToken("1", 1, 4)},
			},
		},
		{name: "{'a b': '1'}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a b", 1, 2),
				token.NewColonToken(1, 3),
				token.NewStringToken("1", 1, 4),
				token.NewRightBraceToken(1, 5),
			},
			expected: []ParsingResult{
				{Path: ".\"a b\"", Token: token.NewStringToken("1", 1, 4)},
			},
		},
		{name: "{'a': '1', 'b': '2', 'c': 3}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewStringToken("1", 1, 4),
				token.NewCommaToken(1, 5),
				token.NewStringToken("b", 1, 6),
				token.NewColonToken(1, 7),
				token.NewStringToken("2", 1, 8),
				token.NewCommaToken(1, 9),
				token.NewStringToken("c", 1, 10),
				token.NewColonToken(1, 111),
				token.NewStringToken("3", 1, 12),
				token.NewRightBraceToken(1, 5),
			},
			expected: []ParsingResult{
				{Path: ".a", Token: token.NewStringToken("1", 1, 4)},
				{Path: ".b", Token: token.NewStringToken("2", 1, 8)},
				{Path: ".c", Token: token.NewStringToken("3", 1, 12)},
			},
		},
		{name: "{'a': { 'b': { 'c': 3}}}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBraceToken(1, 4),
				token.NewStringToken("b", 1, 5),
				token.NewColonToken(1, 6),
				token.NewLeftBraceToken(1, 7),
				token.NewStringToken("c", 1, 8),
				token.NewColonToken(1, 9),
				token.NewStringToken("3", 1, 10),
				token.NewRightBraceToken(1, 11),
				token.NewRightBraceToken(1, 12),
				token.NewRightBraceToken(1, 13),
			},
			expected: []ParsingResult{
				{Path: ".a.b.c", Token: token.NewStringToken("3", 1, 10)},
			},
		},

		// arrays
		{name: "['1', '2']",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
				token.NewStringToken("1", 1, 2),
				token.NewCommaToken(1, 3),
				token.NewStringToken("2", 1, 4),
				token.NewRightBracketToken(1, 5),
			},
			expected: []ParsingResult{
				{Path: ".[0]", Token: token.NewStringToken("1", 1, 2)},
				{Path: ".[1]", Token: token.NewStringToken("2", 1, 4)},
			},
		},
		{name: "['1', ['2']]",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
				token.NewStringToken("1", 1, 2),
				token.NewCommaToken(1, 3),
				token.NewLeftBracketToken(1, 4),
				token.NewStringToken("2", 1, 5),
				token.NewRightBracketToken(1, 6),
				token.NewRightBracketToken(1, 7),
			},
			expected: []ParsingResult{
				{Path: ".[0]", Token: token.NewStringToken("1", 1, 2)},
				{Path: ".[1][0]", Token: token.NewStringToken("2", 1, 5)},
			},
		},
		{name: "[{}]",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
				token.NewLeftBraceToken(1, 3),
				token.NewRightBraceToken(1, 5),
				token.NewRightBracketToken(1, 1),
			},
			expected: []ParsingResult{},
		},
		{name: "[{'a': 1}]",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
				token.NewLeftBraceToken(1, 2),
				token.NewStringToken("a", 1, 3),
				token.NewColonToken(1, 4),
				token.NewNumberToken("1", 1, 5),
				token.NewRightBraceToken(1, 6),
				token.NewRightBracketToken(1, 7),
			},
			expected: []ParsingResult{
				{Path: ".[0].a", Token: token.NewNumberToken("1", 1, 5)},
			},
		},
		{name: "[[[[[[1]]]]]]",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
				token.NewLeftBracketToken(1, 2),
				token.NewLeftBracketToken(1, 3),
				token.NewLeftBracketToken(1, 4),
				token.NewLeftBracketToken(1, 5),
				token.NewLeftBracketToken(1, 6),
				token.NewStringToken("1", 1, 7),
				token.NewRightBracketToken(1, 8),
				token.NewRightBracketToken(1, 9),
				token.NewRightBracketToken(1, 10),
				token.NewRightBracketToken(1, 11),
				token.NewRightBracketToken(1, 12),
				token.NewRightBracketToken(1, 13),
			},
			expected: []ParsingResult{
				{Path: ".[0][0][0][0][0][0]", Token: token.NewStringToken("1", 1, 7)},
			},
		},

		// mixed
		{name: "{'a': ['1', '2']}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBracketToken(1, 4),
				token.NewStringToken("1", 1, 5),
				token.NewCommaToken(1, 6),
				token.NewStringToken("2", 1, 7),
				token.NewRightBracketToken(1, 1),
				token.NewRightBraceToken(1, 5),
			},
			expected: []ParsingResult{
				{Path: ".a[0]", Token: token.NewStringToken("1", 1, 5)},
				{Path: ".a[1]", Token: token.NewStringToken("2", 1, 7)},
			},
		},
		{name: "{'a': [[[[[[1]]]]]]}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBracketToken(1, 4),
				token.NewLeftBracketToken(1, 5),
				token.NewLeftBracketToken(1, 6),
				token.NewLeftBracketToken(1, 7),
				token.NewLeftBracketToken(1, 8),
				token.NewLeftBracketToken(1, 9),
				token.NewStringToken("1", 1, 10),
				token.NewRightBracketToken(1, 11),
				token.NewRightBracketToken(1, 12),
				token.NewRightBracketToken(1, 13),
				token.NewRightBracketToken(1, 14),
				token.NewRightBracketToken(1, 15),
				token.NewRightBracketToken(1, 16),
				token.NewRightBraceToken(1, 17),
			},
			expected: []ParsingResult{
				{Path: ".a[0][0][0][0][0][0]", Token: token.NewStringToken("1", 1, 10)},
			},
		},
		{name: "{'a': ['1', ['2', '3']]}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBracketToken(1, 4),
				token.NewStringToken("1", 1, 5),
				token.NewCommaToken(1, 6),
				token.NewLeftBracketToken(1, 7),
				token.NewStringToken("2", 1, 8),
				token.NewCommaToken(1, 9),
				token.NewStringToken("3", 1, 10),
				token.NewRightBracketToken(1, 11),
				token.NewRightBracketToken(1, 12),
				token.NewRightBraceToken(1, 13),
			},
			expected: []ParsingResult{
				{Path: ".a[0]", Token: token.NewStringToken("1", 1, 5)},
				{Path: ".a[1][0]", Token: token.NewStringToken("2", 1, 8)},
				{Path: ".a[1][1]", Token: token.NewStringToken("3", 1, 10)},
			},
		},
		{name: "{'a': ['1', {'b': ['2', {'c': 3}]}]}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBracketToken(1, 4),
				token.NewStringToken("1", 1, 5),
				token.NewCommaToken(1, 6),
				token.NewLeftBraceToken(1, 7),
				token.NewStringToken("b", 1, 8),
				token.NewColonToken(1, 9),
				token.NewLeftBracketToken(1, 10),
				token.NewStringToken("2", 1, 11),
				token.NewCommaToken(1, 12),
				token.NewLeftBraceToken(1, 13),
				token.NewStringToken("c", 1, 14),
				token.NewColonToken(1, 15),
				token.NewStringToken("3", 1, 16),
				token.NewRightBraceToken(1, 17),
				token.NewRightBracketToken(1, 18),
				token.NewRightBraceToken(1, 19),
				token.NewRightBracketToken(1, 20),
				token.NewRightBraceToken(1, 21),
			},
			expected: []ParsingResult{
				{Path: ".a[0]", Token: token.NewStringToken("1", 1, 5)},
				{Path: ".a[1].b[0]", Token: token.NewStringToken("2", 1, 11)},
				{Path: ".a[1].b[1].c", Token: token.NewStringToken("3", 1, 16)},
			},
		},
		{name: "{'a': ['1', '2'], 'b': '3', 'c': ['4']}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBracketToken(1, 4),
				token.NewStringToken("1", 1, 5),
				token.NewCommaToken(1, 6),
				token.NewStringToken("2", 1, 7),
				token.NewRightBracketToken(1, 8),
				token.NewCommaToken(1, 9),
				token.NewStringToken("b", 1, 10),
				token.NewColonToken(1, 11),
				token.NewStringToken("3", 1, 12),
				token.NewCommaToken(1, 13),
				token.NewStringToken("c", 1, 14),
				token.NewColonToken(1, 15),
				token.NewLeftBracketToken(1, 16),
				token.NewStringToken("4", 1, 17),
				token.NewRightBracketToken(1, 18),
				token.NewRightBraceToken(1, 19),
			},
			expected: []ParsingResult{
				{Path: ".a[0]", Token: token.NewStringToken("1", 1, 5)},
				{Path: ".a[1]", Token: token.NewStringToken("2", 1, 7)},
				{Path: ".b", Token: token.NewStringToken("3", 1, 12)},
				{Path: ".c[0]", Token: token.NewStringToken("4", 1, 17)},
			},
		},
		{name: "{'s': {'t': [[1], -2.0, '3', true, {'x': false}]}}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("s", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBraceToken(1, 4),
				token.NewStringToken("t", 1, 5),
				token.NewColonToken(1, 6),
				token.NewLeftBracketToken(1, 7),
				token.NewLeftBracketToken(1, 8),
				token.NewNumberToken("1", 1, 9),
				token.NewRightBracketToken(1, 10),
				token.NewCommaToken(1, 11),
				token.NewNumberToken("-2.0", 1, 12),
				token.NewCommaToken(1, 13),
				token.NewStringToken("3", 1, 14),
				token.NewCommaToken(1, 15),
				token.NewBooleanToken("true", 1, 16),
				token.NewCommaToken(1, 17),
				token.NewLeftBraceToken(1, 18),
				token.NewStringToken("x", 1, 19),
				token.NewColonToken(1, 20),
				token.NewBooleanToken("false", 1, 21),
				token.NewRightBraceToken(1, 22),
				token.NewRightBracketToken(1, 23),
				token.NewRightBraceToken(1, 24),
				token.NewRightBraceToken(1, 25),
			},
			expected: []ParsingResult{
				{Path: ".s.t[0][0]", Token: token.NewNumberToken("1", 1, 9)},
				{Path: ".s.t[1]", Token: token.NewNumberToken("-2.0", 1, 12)},
				{Path: ".s.t[2]", Token: token.NewStringToken("3", 1, 14)},
				{Path: ".s.t[3]", Token: token.NewBooleanToken("true", 1, 16)},
				{Path: ".s.t[4].x", Token: token.NewBooleanToken("false", 1, 21)},
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
		tokens   []token.Token
		expected string
	}{
		// objects
		{name: "{",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "{{",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewLeftBraceToken(1, 2),
			},
			expected: "invalid JSON. unexpected token { at line 1 column 2"},
		{name: "{{{",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewLeftBraceToken(1, 2),
				token.NewLeftBraceToken(1, 3),
			},
			expected: "invalid JSON. unexpected token { at line 1 column 2"},
		{name: "{{}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewLeftBraceToken(1, 2),
				token.NewRightBraceToken(1, 3),
			},
			expected: "invalid JSON. unexpected token { at line 1 column 2"},
		{name: "{}}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewRightBraceToken(1, 2),
				token.NewRightBraceToken(1, 3),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "{1",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewNumberToken("1", 1, 2),
			},
			expected: "invalid JSON. unexpected token 1 at line 1 column 2"},
		{name: "{true",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewBooleanToken("true", 1, 2),
			},
			expected: "invalid JSON. unexpected token true at line 1 column 2"},
		{name: "{null",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewNullToken(1, 2),
			},
			expected: "invalid JSON. unexpected token null at line 1 column 2"},
		{name: "{:",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewColonToken(1, 2),
			},
			expected: "invalid JSON. unexpected token : at line 1 column 2"},
		{name: "{,",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewCommaToken(1, 2),
			},
			expected: "invalid JSON. unexpected token , at line 1 column 2"},
		{name: "{'a',",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewCommaToken(1, 3),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "{'a',1}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewCommaToken(1, 3),
				token.NewNumberToken("1", 1, 4),
				token.NewRightBraceToken(1, 5),
			},
			expected: "invalid JSON. unexpected token , at line 1 column 3"},
		{name: "{'a','1'}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewCommaToken(1, 3),
				token.NewStringToken("1", 1, 4),
				token.NewRightBraceToken(1, 5),
			},
			expected: "invalid JSON. unexpected token , at line 1 column 3"},
		{name: "{'a': 1,}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewNumberToken("1", 1, 4),
				token.NewCommaToken(1, 5),
				token.NewRightBraceToken(1, 6),
			},
			expected: "invalid JSON. unexpected token } at line 1 column 6"},
		{name: "{'a': 1, 2}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewNumberToken("1", 1, 4),
				token.NewCommaToken(1, 5),
				token.NewNumberToken("2", 1, 6),
				token.NewRightBraceToken(1, 7),
			},
			expected: "invalid JSON. unexpected token 2 at line 1 column 6"},
		{name: "{'a': 'b': 1}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewStringToken("b", 1, 4),
				token.NewColonToken(1, 5),
				token.NewNumberToken("1", 1, 6),
				token.NewRightBraceToken(1, 7),
			},
			expected: "invalid JSON. unexpected token : at line 1 column 5"},
		{name: "{'a': {,}}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBraceToken(1, 4),
				token.NewCommaToken(1, 5),
				token.NewRightBraceToken(1, 6),
				token.NewRightBraceToken(1, 7),
			},
			expected: "invalid JSON. unexpected token , at line 1 column 5"},
		{name: "{'a': {{}}}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBraceToken(1, 4),
				token.NewLeftBraceToken(1, 5),
				token.NewRightBraceToken(1, 6),
				token.NewRightBraceToken(1, 7),
				token.NewRightBraceToken(1, 8),
			},
			expected: "invalid JSON. unexpected token { at line 1 column 5"},

		// arrays
		{name: "[",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "[[",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
				token.NewLeftBracketToken(1, 2),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "[[[",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
				token.NewLeftBracketToken(1, 2),
				token.NewLeftBracketToken(1, 3),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "[[]",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
				token.NewLeftBracketToken(1, 2),
				token.NewRightBracketToken(1, 3),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "[]]",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
				token.NewRightBracketToken(1, 2),
				token.NewRightBracketToken(1, 3),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "[,",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
				token.NewCommaToken(1, 2),
			},
			expected: "invalid JSON. unexpected token , at line 1 column 2"},
		{name: "[1,]",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
				token.NewNumberToken("1", 1, 2),
				token.NewCommaToken(1, 3),
				token.NewRightBracketToken(1, 4),
			},
			expected: "invalid JSON. unexpected token ] at line 1 column 4"},

		// mixed
		{name: "{]",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewRightBracketToken(1, 2),
			},
			expected: "invalid JSON. unexpected token ] at line 1 column 2"},
		{name: "[}",
			tokens: []token.Token{
				token.NewLeftBracketToken(1, 1),
				token.NewRightBraceToken(1, 2),
			},
			expected: "invalid JSON. unexpected token } at line 1 column 2"},
		{name: "{[}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewLeftBracketToken(1, 2),
				token.NewRightBraceToken(1, 3),
			},
			expected: "invalid JSON. unexpected token [ at line 1 column 2"},
		{name: "{'a': [}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBracketToken(1, 4),
				token.NewRightBraceToken(1, 5),
			},
			expected: "invalid JSON. unexpected token } at line 1 column 5"},
		{name: "{'a': [,",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBracketToken(1, 4),
				token.NewCommaToken(1, 5),
			},
			expected: "invalid JSON. unexpected token , at line 1 column 5"},
		{name: "{'a': [,",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBracketToken(1, 4),
				token.NewCommaToken(1, 5),
			},
			expected: "invalid JSON. unexpected token , at line 1 column 5"},
		{name: "{'a': [{},",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("a", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBracketToken(1, 4),
				token.NewLeftBraceToken(1, 5),
				token.NewRightBraceToken(1, 6),
			},
			expected: "invalid JSON. Unexpected end of JSON input"},
		{name: "{'s': {'t': [[1], -2.0, '3', true, {'x': false,}]}}",
			tokens: []token.Token{
				token.NewLeftBraceToken(1, 1),
				token.NewStringToken("s", 1, 2),
				token.NewColonToken(1, 3),
				token.NewLeftBraceToken(1, 4),
				token.NewStringToken("t", 1, 5),
				token.NewColonToken(1, 6),
				token.NewLeftBracketToken(1, 7),
				token.NewLeftBracketToken(1, 8),
				token.NewNumberToken("1", 1, 9),
				token.NewRightBracketToken(1, 10),
				token.NewCommaToken(1, 11),
				token.NewNumberToken("-2.0", 1, 12),
				token.NewCommaToken(1, 13),
				token.NewStringToken("3", 1, 14),
				token.NewCommaToken(1, 15),
				token.NewBooleanToken("true", 1, 16),
				token.NewCommaToken(1, 17),
				token.NewLeftBraceToken(1, 18),
				token.NewStringToken("x", 1, 19),
				token.NewColonToken(1, 20),
				token.NewBooleanToken("false", 1, 21),
				token.NewCommaToken(1, 22),
				token.NewRightBraceToken(1, 23),
				token.NewRightBracketToken(1, 24),
				token.NewRightBraceToken(1, 25),
				token.NewRightBraceToken(1, 26),
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
