package parser

import "testing"

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
