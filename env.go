package calculon

type Function = func(args []float64) (float64, error)

type Env struct {
	vars  map[string]float64
	funcs map[string]Function
}

func NewEnv() *Env {
	return &Env{
		vars:  make(map[string]float64),
		funcs: make(map[string]Function),
	}
}

func (e *Env) SetVar(name string, value float64) {
	e.vars[name] = value
}

func (e *Env) SetFunc(name string, fn Function) {
	e.funcs[name] = fn
}

func (e *Env) LookupVar(name string) (float64, bool) {
	val, found := e.vars[name]
	return val, found
}

func (e *Env) LookupFunc(name string) (Function, bool) {
	fn, found := e.funcs[name]
	return fn, found
}

func MathEnv() *Env {
	e := NewEnv()
	for name, val := range builtinVars {
		e.vars[name] = val
	}

	for name, fn := range builtinFuncs {
		e.funcs[name] = fn
	}

	return e
}
