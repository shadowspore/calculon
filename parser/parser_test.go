package parser

import (
	"fmt"
	"testing"

	. "github.com/zweihander/calculon/ast"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Node
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
								Inner: BinaryOp{
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
								Inner: BinaryOp{
									Op: "+",
									Left: UnaryOp{
										Op:      "-",
										Operand: Number{Value: 9},
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
					Inner: BinaryOp{
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
							Args: []Node{
								BinaryOp{
									Op:    "+",
									Left:  UnaryOp{Op: "-", Operand: Number{Value: 2}},
									Right: Variable{Name: "y"},
								},
							},
						},
						Right: Number{Value: 3},
					},
				},
				Right: FunctionCall{
					Name: "baz",
					Args: []Node{
						Variable{Name: "z"},
						FunctionCall{
							Name: "bax",
							Args: []Node{Variable{Name: "x"}},
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
			require.Equal(t, test.err, err)

			require.Equal(t, test.expected, expr)
		})
	}
}
