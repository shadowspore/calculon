package calculon

import (
	"fmt"
	"math"
)

type Function = func(args []float64) (float64, error)

var (
	builtinVars = map[string]float64{
		"Pi": math.Pi,
		"E":  math.E,
	}

	builtinFuncs = map[string]Function{
		"sin": func(args []float64) (float64, error) {
			if len(args) != 1 {
				return 0, fmt.Errorf("sin() requires 1 arg")
			}

			return math.Sin(args[0]), nil
		},
		"cos": func(args []float64) (float64, error) {
			if len(args) != 1 {
				return 0, fmt.Errorf("cos() requires 1 arg")
			}

			return math.Cos(args[0]), nil
		},
	}
)
