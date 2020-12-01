package calculon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
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

	for _, test := range tests {
		expr, err := Parse(test.input)
		assert.NoError(t, err)

		result, err := expr.Eval(EmptyContext{})
		assert.NoError(t, err)

		assert.Equal(t, test.expected, result)
	}
}
