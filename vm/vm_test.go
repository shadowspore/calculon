package vm_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zweihander/calculon/compiler"
	"github.com/zweihander/calculon/parser"
	"github.com/zweihander/calculon/vm"
)

func TestVM(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"2+2", 4},
		{"2*4", 8},
		{"2*(3-4)+2/4", -1.5},
		{"2 * (3 + 4) / 5 - 512 * (-9 + 10) * 332 - 55 / 2", -170008.7},
		{"4^3^2", 262144},
	}

	vm := vm.New(vm.Config{})
	for _, test := range tests {
		expr, err := parser.Parse(test.input)
		require.NoError(t, err)

		program, err := compiler.Compile(expr)
		require.NoError(t, err)

		result, err := vm.Run(program, nil)
		require.NoError(t, err)

		require.Equal(t, test.expected, result)
	}
}
