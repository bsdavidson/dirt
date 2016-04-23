package parser

import (
	"testing"
)

func TestParse(t *testing.T) {
	type testCase struct {
		input  string
		tokens []Token
		err    error
	}
	testCases := []testCase{
		testCase{
			input: "this is a test",
			tokens: []Token{
				Token{Value: "this", Start: 0, Width: 4},
				Token{Value: "is", Start: 5, Width: 2},
				Token{Value: "a", Start: 8, Width: 1},
				Token{Value: "test", Start: 10, Width: 4},
			},
			err: nil,
		},
		testCase{
			input: "  this   has \t  weird whitespace\r\n   ",
			tokens: []Token{
				Token{Value: "this", Start: 2, Width: 4},
				Token{Value: "has", Start: 9, Width: 3},
				Token{Value: "weird", Start: 16, Width: 5},
				Token{Value: "whitespace", Start: 22, Width: 10},
			},
		},
		testCase{
			input: `this has 'single quoted' and "double quoted with \"escaped\"" characters`,
			tokens: []Token{
				Token{Value: "this", Start: 0, Width: 4},
				Token{Value: "has", Start: 5, Width: 3},
				Token{Value: "single quoted", Start: 9, Width: 15},
				Token{Value: "and", Start: 25, Width: 3},
				Token{Value: `double quoted with "escaped"`, Start: 29, Width: 32},
				Token{Value: "characters", Start: 62, Width: 10},
			},
		},
	}

	for i, tc := range testCases {
		tokens, err := Parse(tc.input)
		if err != tc.err {
			t.Errorf("case %d: expected err to be %+v, was %+v", i, tc.err, err)
		} else if len(tokens) != len(tc.tokens) {
			t.Errorf("case %d: expected tokens to be %+v, was %+v", i, tc.tokens, tokens)
		} else {
			for ti, token := range tc.tokens {
				if tokens[ti] != token {
					t.Errorf("case %d: expected tokens[%d] to be %+v, was %+v", i, ti, token, tokens[ti])
					break
				}
			}
		}
	}
}
