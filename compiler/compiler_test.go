package compiler_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zweihander/calculon/ast"
	"github.com/zweihander/calculon/compiler"
	"github.com/zweihander/calculon/parser"
	"github.com/zweihander/calculon/vm"
)

func TestCompiler(t *testing.T) {
	tests := []struct {
		input    ast.Node
		expected vm.Program
	}{
		{
			input: mustParse("2 + 2"),
			expected: vm.Program{
				Opcodes: []vm.Opcode{
					vm.Num, vm.Num, vm.Add,
				},
				Values: []float64{
					2, 2,
				},
			},
		},
		{
			input: mustParse("1 * 3 + 2 / 4"),
			expected: vm.Program{
				Opcodes: []vm.Opcode{
					vm.Num, vm.Num, vm.Mult,
					vm.Num, vm.Num, vm.Div,
					vm.Add,
				},
				Values: []float64{
					1, 3, 2, 4,
				},
			},
		},
	}

	for _, test := range tests {
		program, err := compiler.Compile(test.input)
		require.NoError(t, err)

		require.Equal(t, test.expected, program)
	}
}

func mustParse(input string) ast.Node {
	n, err := parser.Parse(input)
	if err != nil {
		panic(err)
	}

	return n
}
