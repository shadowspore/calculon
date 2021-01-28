package compiler_test

import (
	"testing"

	"github.com/zweihander/calculon/compiler"
	"github.com/zweihander/calculon/parser"
)

func BenchmarkCompiler(b *testing.B) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "simple",
			input: "2 + 2",
		},
		{
			name:  "complex",
			input: "2 * (3 + 4) / 1024 - 512 * (-9 + 100) * 1533223 - 55 / 2",
		},
		{
			name:  "complex-with-vars",
			input: "2 * (x + 4) / y - 512 * (9 + 100) * z - 55 / 2",
		},
		{
			name:  "simple-sin-func",
			input: "sin(5)",
		},
	}

	b.ReportAllocs()
	for _, test := range tests {
		node, err := parser.Parse(test.input)
		if err != nil {
			b.Error(err)
		}

		b.Run(test.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := compiler.Compile(node)
				if err != nil {
					b.Error(err)
				}
			}
		})
	}
}
