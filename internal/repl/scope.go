package repl

import "github.com/xjem/calculon"

var _ calculon.EvalContext = (*Scope)(nil)

type Scope struct {
	parent calculon.EvalContext
	vars   map[string]float64
	funcs  map[string]calculon.Function
}

func (s *Scope) SetVar(name string, value float64) { s.vars[name] = value }

func (s *Scope) SetFunc(name string, fn calculon.Function) { s.funcs[name] = fn }

func (s *Scope) LookupVar(name string) (float64, bool) {
	val, found := s.vars[name]
	if found {
		return val, true
	}

	return s.parent.LookupVar(name)
}

func (s *Scope) LookupFunc(name string) (calculon.Function, bool) {
	fn, found := s.funcs[name]
	if found {
		return fn, true
	}

	return s.parent.LookupFunc(name)
}

func NewScope(parent calculon.EvalContext) *Scope {
	if parent == nil {
		parent = calculon.EmptyContext{}
	}
	return &Scope{
		parent: parent,
		vars:   map[string]float64{},
		funcs:  map[string]calculon.Function{},
	}
}
