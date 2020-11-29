package calculon

import (
	"testing"
)

func BenchmarkParser(b *testing.B) {
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
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := Parse(test.input)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkEvaler(b *testing.B) {
	tests := []struct {
		name  string
		input string
		ctx   EvalContext
	}{
		{
			name:  "simple",
			input: "2 + 2",
			ctx:   EmptyContext{},
		},
		{
			name:  "complex",
			input: "2 * (3 + 4) / 1024 - 512 * (-9 + 100) * 1533223 - 55 / 2",
			ctx:   EmptyContext{},
		},
		{
			name:  "complex-with-vars",
			input: "2 * (x + 4) / y - 512 * (9 + 100) * z - 55 / 2",
			ctx: createCtx(map[string]float64{
				"x": 1,
				"y": 2,
				"z": 3,
			}, nil),
		},
		{
			name:  "simple-sin-func",
			input: "sin(5)",
			ctx:   MathContext(),
		},
	}

	for _, test := range tests {
		expr, err := Parse(test.input)
		if err != nil {
			b.Fatal(err)
		}

		b.Run(test.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := expr.Eval(test.ctx)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func createCtx(vars map[string]float64, funcs map[string]Function) EvalContext {
	ctx := NewContext()
	for name, val := range vars {
		ctx.SetVar(name, val)
	}

	for name, fn := range funcs {
		ctx.SetFunc(name, fn)
	}

	return ctx
}
