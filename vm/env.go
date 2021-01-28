package vm

type Environment interface {
	LookupVar(name string) (float64, bool)
	LookupFunc(name string) (Function, bool)
}

type Function = func(args []float64) (float64, error)
