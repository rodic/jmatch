package tokenizer

// import (
// 	"reflect"
// 	"testing"

// 	token "github.com/rodic/jmatch/token"
// )

// func TestTokenizeValidInputs(t *testing.T) {
// 	testCases := []struct {
// 		name     string
// 		input    string
// 		expected []token.Token
// 	}{
// 		// Simple cases, {} and one key value pair
// 		{name: "empty",
// 			input: "{}",
// 			expected: []token.Token{
// 				token.NewLeftBraceToken(1, 1),
// 				token.NewRightBraceToken(1, 2)}},
// 		{name: "emptyWithSpaces",
// 			input: "{     }",
// 			expected: []token.Token{
// 				token.NewLeftBraceToken(1, 1),
// 				token.NewRightBraceToken(1, 7)}},
// 		{name: "simplePair",
// 			input: "{\"a\":\"1\"}",
// 			expected: []token.Token{
// 				token.NewLeftBraceToken(1, 1),
// 				token.NewStringToken("a", 1, 2),
// 				token.NewColonToken(1, 5),
// 				token.NewStringToken("1", 1, 6),
// 				token.NewRightBraceToken(1, 9)}},
// 		{name: "simplePairWithSpaces",
// 			input: "{  \"a\"  :   \"1\"  }",
// 			expected: []token.Token{
// 				token.NewLeftBraceToken(1, 1),
// 				token.NewStringToken("a", 1, 4),
// 				token.NewColonToken(1, 9),
// 				token.NewStringToken("1", 1, 13),
// 				token.NewRightBraceToken(1, 18)}},
// 		{name: "simplePairWithInt",
// 			input: "{\"a\":1}",
// 			expected: []token.Token{
// 				token.NewLeftBraceToken(1, 1),
// 				token.NewStringToken("a", 1, 2),
// 				token.NewColonToken(1, 5),
// 				token.NewNumberToken("1", 1, 6),
// 				token.NewRightBraceToken(1, 7)}},
// 		{name: "simplePairWithNegativeInt",
// 			input: "{\"a\":-1}",
// 			expected: []token.Token{
// 				token.NewLeftBraceToken(1, 1),
// 				token.NewStringToken("a", 1, 2),
// 				token.NewColonToken(1, 5),
// 				token.NewNumberToken("-1", 1, 6),
// 				token.NewRightBraceToken(1, 8)}},
// 		{name: "simplePairWithIntAndSpaces",
// 			input: "{  \"a\"  : 1  }",
// 			expected: []token.Token{
// 				token.NewLeftBraceToken(1, 1),
// 				token.NewStringToken("a", 1, 4),
// 				token.NewColonToken(1, 9),
// 				token.NewNumberToken("1", 1, 11),
// 				token.NewRightBraceToken(1, 14)}},
// 		{name: "simplePairWithDecimal",
// 			input: "{\"a\":1.01}",
// 			expected: []token.Token{
// 				token.NewLeftBraceToken(1, 1),
// 				token.NewStringToken("a", 1, 2),
// 				token.NewColonToken(1, 5),
// 				token.NewNumberToken("1.01", 1, 6),
// 				token.NewRightBraceToken(1, 10)}},
// 		{name: "simplePairWithTrue",
// 			input: "{\"a\":true}",
// 			expected: []token.Token{
// 				token.NewLeftBraceToken(1, 1),
// 				token.NewStringToken("a", 1, 2),
// 				token.NewColonToken(1, 5),
// 				token.NewBooleanToken("true", 1, 6),
// 				token.NewRightBraceToken(1, 10)}},
// 		{name: "simplePairWithFalse",
// 			input: "{\"a\":false}",
// 			expected: []token.Token{
// 				token.NewLeftBraceToken(1, 1),
// 				token.NewStringToken("a", 1, 2),
// 				token.NewColonToken(1, 5),
// 				token.NewBooleanToken("false", 1, 6),
// 				token.NewRightBraceToken(1, 11)}},
// 		{name: "simplePairWithNull",
// 			input: "{\"a\":null}",
// 			expected: []token.Token{
// 				token.NewLeftBraceToken(1, 1),
// 				token.NewStringToken("a", 1, 2),
// 				token.NewColonToken(1, 5),
// 				token.NewNullToken(1, 6),
// 				token.NewRightBraceToken(1, 10)}},
// 		{name: "simplePairWithDotInKey",
// 			input: "{\"a.b\":\"1\"}",
// 			expected: []token.Token{
// 				token.NewLeftBraceToken(1, 1),
// 				token.NewStringToken("a.b", 1, 2),
// 				token.NewColonToken(1, 7),
// 				token.NewStringToken("1", 1, 8),
// 				token.NewRightBraceToken(1, 11)}},
// 		{name: "simplePairWithSpaceInKey",
// 			input: "{\"a b\":\"1\"}",
// 			expected: []token.Token{
// 				token.NewLeftBraceToken(1, 1),
// 				token.NewStringToken("a b", 1, 2),
// 				token.NewColonToken(1, 7),
// 				token.NewStringToken("1", 1, 8),
// 				token.NewRightBraceToken(1, 11)}},

// 		// Array
// 		{name: "simplePairWithArray",
// 			input: "{\"a\":[1, \"2\", true, null]}",
// 			expected: []token.Token{
// 				token.NewLeftBraceToken(1, 1),
// 				token.NewStringToken("a", 1, 2),
// 				token.NewColonToken(1, 5),
// 				token.NewLeftBracketToken(1, 6),
// 				token.NewNumberToken("1", 1, 7),
// 				token.NewCommaToken(1, 8),
// 				token.NewStringToken("2", 1, 10),
// 				token.NewCommaToken(1, 13),
// 				token.NewBooleanToken("true", 1, 15),
// 				token.NewCommaToken(1, 19),
// 				token.NewNullToken(1, 21),
// 				token.NewRightBracketToken(1, 25),
// 				token.NewRightBraceToken(1, 26)}},

// 		// Nested
// 		{name: "nested",
// 			input: "{\"a\": \"ünicode\", \"b\"  : { \"list\": [true, [false, null, {  }]]}}",
// 			expected: []token.Token{
// 				token.NewLeftBraceToken(1, 1),
// 				token.NewStringToken("a", 1, 2),
// 				token.NewColonToken(1, 5),
// 				token.NewStringToken("ünicode", 1, 7),

// 				token.NewCommaToken(1, 16),
// 				token.NewStringToken("b", 1, 18),
// 				token.NewColonToken(1, 23),
// 				token.NewLeftBraceToken(1, 25),
// 				token.NewStringToken("list", 1, 27),

// 				token.NewColonToken(1, 33),
// 				token.NewLeftBracketToken(1, 35),
// 				token.NewBooleanToken("true", 1, 36),
// 				token.NewCommaToken(1, 40),
// 				token.NewLeftBracketToken(1, 42),
// 				token.NewBooleanToken("false", 1, 43),
// 				token.NewCommaToken(1, 48),
// 				token.NewNullToken(1, 50),
// 				token.NewCommaToken(1, 54),
// 				token.NewLeftBraceToken(1, 56),
// 				token.NewRightBraceToken(1, 59),
// 				token.NewRightBracketToken(1, 60),
// 				token.NewRightBracketToken(1, 61),
// 				token.NewRightBraceToken(1, 62),
// 				token.NewRightBraceToken(1, 63),
// 			}},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			tokenizer := NewTokenizer(tc.input)
// 			result, err := tokenizer.Tokenize()
// 			if err != nil {
// 				t.Error(err)
// 			}
// 			if !reflect.DeepEqual(result, tc.expected) {
// 				t.Errorf("Expected '%v', got '%v' instead\n", tc.expected, result)
// 			}
// 		})
// 	}
// }

// func TestTokenizeInvalidInputs(t *testing.T) {
// 	testCases := []struct {
// 		name     string
// 		input    string
// 		expected string
// 	}{
// 		{name: "invalidNumber",
// 			input:    "{\"a\":1.2.3}",
// 			expected: "invalid JSON. unexpected token . at line 1 column 9"},
// 		{name: "invalidText",
// 			input:    "{\"a\":    truef}",
// 			expected: "invalid JSON. unexpected token truef at line 1 column 10"},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			tokenizer := NewTokenizer(tc.input)
// 			res, err := tokenizer.Tokenize()
// 			if err == nil {
// 				t.Errorf("Expected error %s got result %v", tc.expected, res)
// 			}
// 			if err.Error() != tc.expected {
// 				t.Errorf("Expected error %s got %s", tc.expected, err)
// 			}

// 		})
// 	}
// }
