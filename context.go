package calculon

type EvalContext interface {
	LookupVar(name string) (float64, bool)
	LookupFunc(name string) (Function, bool)
}

type EmptyContext struct{}

func (EmptyContext) LookupVar(name string) (float64, bool) { return 0, false }

func (EmptyContext) LookupFunc(name string) (Function, bool) { return nil, false }

// Context contains user-defined variables and functions.
type Context struct {
	vars  map[string]float64
	funcs map[string]Function
}

func NewContext() *Context {
	return &Context{
		vars:  make(map[string]float64),
		funcs: make(map[string]Function),
	}
}

func (ctx *Context) SetVar(name string, value float64) {
	ctx.vars[name] = value
}

func (ctx *Context) SetFunc(name string, fn Function) {
	ctx.funcs[name] = fn
}

func (ctx *Context) ForEachVars(iter func(name string, value float64)) {
	for name, val := range ctx.vars {
		iter(name, val)
	}
}

func (ctx *Context) ForEachFuncs(iter func(name string, fn Function)) {
	for name, fn := range ctx.funcs {
		iter(name, fn)
	}
}

func (ctx *Context) LookupVar(name string) (float64, bool) {
	val, found := ctx.vars[name]
	return val, found
}

func (ctx *Context) LookupFunc(name string) (Function, bool) {
	fn, found := ctx.funcs[name]
	return fn, found
}

func MathContext() *Context {
	ctx := NewContext()
	for name, val := range builtinVars {
		ctx.vars[name] = val
	}

	for name, fn := range builtinFuncs {
		ctx.funcs[name] = fn
	}

	return ctx
}
