package vm_test

import (
	"math"
	"testing"

	"github.com/zweihander/calculon"
	"github.com/zweihander/calculon/compiler"
	"github.com/zweihander/calculon/parser"
	"github.com/zweihander/calculon/vm"
)

func BenchmarkVM(b *testing.B) {
	tests := []struct {
		name  string
		input string
		env   vm.Environment
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
			env: createCtx(map[string]float64{
				"x": 1,
				"y": 2,
				"z": 3,
			}, nil),
		},
		{
			name:  "simple-sin-func",
			input: "sin(5)",
			env: createCtx(nil, map[string]func(args []float64) (float64, error){
				"sin": func(args []float64) (float64, error) {
					return math.Sin(args[0]), nil
				},
			}),
		},
	}

	vm := vm.New(vm.Config{})
	b.ReportAllocs()
	for _, test := range tests {
		expr, err := parser.Parse(test.input)
		if err != nil {
			b.Fatal(err)
		}

		program, err := compiler.Compile(expr)
		if err != nil {
			b.Fatal(err)
		}

		b.Run(test.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := vm.Run(program, test.env)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func createCtx(vars map[string]float64, funcs map[string]vm.Function) vm.Environment {
	e := calculon.NewEnv()
	for name, val := range vars {
		e.SetVar(name, val)
	}

	for name, fn := range funcs {
		e.SetFunc(name, fn)
	}

	return e
}
