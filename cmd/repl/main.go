package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/zweihander/calculon"
)

// for functions, to keep the global scope access
type ForkCtx struct {
	parent calculon.EvalContext
	*calculon.Context
}

func (fc *ForkCtx) LookupVar(name string) (float64, bool) {
	val, found := fc.Context.LookupVar(name)
	if found {
		return val, true
	}

	return fc.parent.LookupVar(name)
}

func (fc *ForkCtx) LookupFunc(name string) (calculon.Function, bool) {
	fn, found := fc.Context.LookupFunc(name)
	if found {
		return fn, true
	}

	return fc.parent.LookupFunc(name)
}

func NewForkCtx(parent calculon.EvalContext) *ForkCtx {
	return &ForkCtx{
		parent:  parent,
		Context: calculon.NewContext(),
	}
}

func main() {
	r := bufio.NewReader(os.Stdin)
	ctx := calculon.MathContext()
	for {
		fmt.Print(">> ")
		input, err := r.ReadString('\n')
		if err != nil {
			panic(err)
		}

		if err := exec(strings.TrimSpace(input), ctx); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			continue
		}
	}
}

func exec(input string, ctx *calculon.Context) error {
	if strings.Contains(input, "=") {
		return define(input, ctx)
	}

	expr, err := calculon.Parse(input)
	if err != nil {
		return fmt.Errorf("parse expression: %w", err)
	}

	result, err := expr.Eval(ctx)
	if err != nil {
		return fmt.Errorf("eval: %w", err)
	}

	fmt.Println(result)
	return nil
}

func define(input string, ctx *calculon.Context) error {
	vars := strings.Split(input, "=")
	if len(vars) != 2 {
		return fmt.Errorf("bad def")
	}

	what, err := calculon.Parse(strings.TrimSpace(vars[0]))
	if err != nil {
		return err
	}

	def, err := calculon.Parse(strings.TrimSpace(vars[1]))
	if err != nil {
		return err
	}

	switch w := what.(type) {
	case calculon.Variable:
		num, ok := def.(calculon.Number)
		if !ok {
			return fmt.Errorf("invalid value type: %T", def)
		}

		ctx.SetVar(w.Name, num.Value)
	case calculon.FunctionCall:
		var pnames []string
		for _, arg := range w.Args {
			vararg, ok := arg.(calculon.Variable)
			if !ok {
				return fmt.Errorf("unsupported function argument: %T", arg)
			}

			pnames = append(pnames, vararg.Name)
		}

		fnCtx := NewForkCtx(ctx)
		ctx.SetFunc(w.Name, func(args []float64) (float64, error) {
			if len(args) != len(pnames) {
				return 0, fmt.Errorf("%s(): bad params count (want %d, got %d)", w.Name, len(pnames), len(args))
			}

			for i, paramName := range pnames {
				fnCtx.SetVar(paramName, args[i])
			}

			return def.Eval(fnCtx)
		})
	default:
		return fmt.Errorf("bad def type: %T", w)
	}

	return nil
}