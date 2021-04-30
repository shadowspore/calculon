package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "simple",
			input: "-(2.53+3)",
			expected: []Token{
				{Minus, ""},
				{OpenParen, ""},
				{Number, "2.53"},
				{Plus, ""},
				{Number, "3"},
				{CloseParen, ""},
			},
		},
		{
			name:  "simple2",
			input: "(x)",
			expected: []Token{
				{OpenParen, ""},
				{Ident, "x"},
				{CloseParen, ""},
			},
		},
		{
			name:  "medium",
			input: " sin(5) + foo(Pi, 3)^4 ",
			expected: []Token{
				{Ident, "sin"},
				{OpenParen, ""},
				{Number, "5"},
				{CloseParen, ""},
				{Plus, ""},
				{Ident, "foo"},
				{OpenParen, ""},
				{Ident, "Pi"},
				{Comma, ""},
				{Number, "3"},
				{CloseParen, ""},
				{Caret, ""},
				{Number, "4"},
			},
		},
		{
			name:     "empty",
			input:    "",
			expected: nil,
		},
	}

	for _, test := range tests {
		l := New(test.input)

		var tokens []Token
		for tok := l.Next(); tok.Kind != EOF; tok = l.Next() {
			tokens = append(tokens, tok)
		}

		assert.Equal(t, test.expected, tokens)
	}
}
