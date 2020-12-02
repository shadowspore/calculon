package calculon

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Expression
		err      error
	}{
		{
			name:  "simple",
			input: "2 + 2",
			expected: BinaryOp{
				Op:    "+",
				Left:  Number{Value: 2},
				Right: Number{Value: 2},
			},
		},
		{
			name:  "complex",
			input: "2 * (3 + 4) / 1024 - 512 * (-9 + 100) * 1533223 - 55 / 2",
			expected: BinaryOp{
				Op: "-",
				Left: BinaryOp{
					Op: "-",
					Left: BinaryOp{
						Op: "/",
						Left: BinaryOp{
							Op:   "*",
							Left: Number{Value: 2},
							Right: Parentheses{
								Expr: BinaryOp{
									Op:    "+",
									Left:  Number{Value: 3},
									Right: Number{Value: 4},
								},
							},
						},
						Right: Number{Value: 1024},
					},
					Right: BinaryOp{
						Op: "*",
						Left: BinaryOp{
							Op:   "*",
							Left: Number{Value: 512},
							Right: Parentheses{
								Expr: BinaryOp{
									Op: "+",
									Left: UnaryOp{
										Op:   "-",
										Expr: Number{Value: 9},
									},
									Right: Number{Value: 100},
								},
							},
						},
						Right: Number{Value: 1533223},
					},
				},
				Right: BinaryOp{
					Op:    "/",
					Left:  Number{Value: 55},
					Right: Number{Value: 2},
				},
			},
		},
		{
			name:  "simple-with-vars",
			input: "x + (y / 3)",
			expected: BinaryOp{
				Op:   "+",
				Left: Variable{Name: "x"},
				Right: Parentheses{
					Expr: BinaryOp{
						Op:    "/",
						Left:  Variable{Name: "y"},
						Right: Number{Value: 3},
					},
				},
			},
		},
		{
			name:  "simple-with-functions",
			input: "foo() + bar(-2 + y)^3 - baz(z, bax(x))",
			expected: BinaryOp{
				Op: "-",
				Left: BinaryOp{
					Op:   "+",
					Left: FunctionCall{Name: "foo"},
					Right: BinaryOp{
						Op: "^",
						Left: FunctionCall{
							Name: "bar",
							Args: []Expression{
								BinaryOp{
									Op:    "+",
									Left:  UnaryOp{Op: "-", Expr: Number{Value: 2}},
									Right: Variable{Name: "y"},
								},
							},
						},
						Right: Number{Value: 3},
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
			expected: BinaryOp{
				Op:   "^",
				Left: Number{Value: 2},
				Right: BinaryOp{
					Op:    "^",
					Left:  Number{Value: 3},
					Right: Number{Value: 4},
				},
			},
		},
		{
			name:  "simple-modulo",
			input: "2 + 3 % 4",
			expected: BinaryOp{
				Op:   "+",
				Left: Number{Value: 2},
				Right: BinaryOp{
					Op:    "%",
					Left:  Number{Value: 3},
					Right: Number{Value: 4},
				},
			},
		},
		{
			name:  "wrong1",
			input: "f)(",
			err:   fmt.Errorf("unexpected token: Kind:)"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr, err := Parse(test.input)
			assert.Equal(t, test.err, err)

			assert.Equal(t, test.expected, expr)
		})
	}
}
