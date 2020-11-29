package calculon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Expression
	}{
		{
			name:  "simple",
			input: "2 + 2",
			expected: Addition{
				Left:  Number{Value: 2},
				Right: Number{Value: 2},
			},
		},
		{
			name:  "complex",
			input: "2 * (3 + 4) / 1024 - 512 * (-9 + 100) * 1533223 - 55 / 2",
			expected: Subtraction{
				Left: Subtraction{
					Left: Division{
						Left: Multiplication{
							Left: Number{Value: 2},
							Right: Parentheses{
								Expr: Addition{
									Left:  Number{Value: 3},
									Right: Number{Value: 4},
								},
							},
						},
						Right: Number{Value: 1024},
					},
					Right: Multiplication{
						Left: Multiplication{
							Left: Number{Value: 512},
							Right: Parentheses{
								Expr: Addition{
									Left: Negation{
										Expr: Number{Value: 9},
									},
									Right: Number{Value: 100},
								},
							},
						},
						Right: Number{Value: 1533223},
					},
				},
				Right: Division{
					Left:  Number{Value: 55},
					Right: Number{Value: 2},
				},
			},
		},
		{
			name:  "simple-with-vars",
			input: "x + (y / 3)",
			expected: Addition{
				Left: Variable{Name: "x"},
				Right: Parentheses{
					Expr: Division{
						Left:  Variable{Name: "y"},
						Right: Number{Value: 3},
					},
				},
			},
		},
		{
			name:  "simple-with-functions",
			input: "foo() + bar(-2 + y)^3 - baz(z, bax(x))",
			expected: Subtraction{
				Left: Addition{
					Left: FunctionCall{Name: "foo"},
					Right: Exponentiation{
						Num: FunctionCall{
							Name: "bar",
							Args: []Expression{
								Addition{
									Left:  Negation{Expr: Number{Value: 2}},
									Right: Variable{Name: "y"},
								},
							},
						},
						Power: Number{Value: 3},
					},
				},
				Right: FunctionCall{
					Name: "baz",
					Args: []Expression{
						Variable{Name: "z"},
						FunctionCall{
							Name: "bax",
							Args: []Expression{Variable{Name: "x"}},
						},
					},
				},
			},
		},
		{
			name:  "multi-exponent",
			input: "2^3^4",
			expected: Exponentiation{
				Num: Number{Value: 2},
				Power: Exponentiation{
					Num:   Number{Value: 3},
					Power: Number{Value: 4},
				},
			},
		},
		{
			name:  "simple-modulo",
			input: "2 + 3 % 4",
			expected: Addition{
				Left: Number{Value: 2},
				Right: Modulo{
					Left:  Number{Value: 3},
					Right: Number{Value: 4},
				},
			},
		},
		{
			name:  "wrong1",
			input: "f)(",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr, err := Parse(test.input)
			assert.NoError(t, err)

			assert.Equal(t, test.expected, expr)
		})
	}
}
