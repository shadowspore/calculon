package vm

import (
	"fmt"
	"math"
)

type VM struct {
	pool *f64pool
}

func New(cfg Config) *VM {
	cfg.setDefaults()

	return &VM{
		pool: &f64pool{
			max: cfg.StackPoolSize,
			cap: cfg.StackCapacity,
		},
	}
}

func (vm *VM) Run(p Program, env Environment) (float64, error) {
	var (
		stack      = vm.pool.Get()
		valOffset  int
		varOffset  int
		funcOffset int
	)

	defer func() { stack = stack[:0]; vm.pool.Put(stack) }()

	pop := func() (float64, bool) {
		if len(stack) == 0 {
			return 0, false
		}

		val := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		return val, true
	}

	popN := func(n int) ([]float64, bool) {
		if len(stack) < n {
			return nil, false
		}

		vals := stack[len(stack)-n:]
		stack = stack[:len(stack)-n]
		return vals, true
	}

	push := func(val float64) { stack = append(stack, val) }

	for _, op := range p.Opcodes {
		switch op {
		case Num:
			push(p.Values[valOffset])
			valOffset++
		case Add, Sub, Mult, Div, Mod, Exp:
			r, ok := pop()
			if !ok {
				return 0, fmt.Errorf("invalid program")
			}

			l, ok := pop()
			if !ok {
				return 0, fmt.Errorf("invalid program")
			}

			switch op {
			case Add:
				push(l + r)
			case Sub:
				push(l - r)
			case Mult:
				push(l * r)
			case Div:
				if r == 0 {
					return 0, fmt.Errorf("divide by zero")
				}

				push(l / r)
			case Mod:
				push(math.Mod(l, r))
			case Exp:
				push(math.Pow(l, r))
			}
		case Neg:
			val, ok := pop()
			if !ok {
				return 0, fmt.Errorf("invalid program")
			}

			push(-val)
		case Var:
			varName := p.Vars[varOffset]
			varOffset++

			val, ok := env.LookupVar(varName)
			if !ok {
				return 0, fmt.Errorf("variable not set: %s", varName)
			}

			push(val)
		case Func:
			meta := p.Funcs[funcOffset]
			funcOffset++

			if len(p.Values) < meta.ArgsCount {
				return 0, fmt.Errorf("invalid program")
			}

			f, ok := env.LookupFunc(meta.Name)
			if !ok {
				return 0, fmt.Errorf("function not defined: %s", meta.Name)
			}

			args, ok := popN(meta.ArgsCount)
			if !ok {
				return 0, fmt.Errorf("invalid args count")
			}

			result, err := f(args)
			if err != nil {
				return 0, err
			}

			push(result)
		default:
			return 0, fmt.Errorf("unknown op: %v", op)
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid program stack")
	}

	return stack[0], nil
}
