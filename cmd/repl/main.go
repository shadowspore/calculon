package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/xjem/calculon"
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
		if err := define(input, ctx); err != nil {
			return fmt.Errorf("assign: %w", err)
		}

		return nil
	}

	switch strings.ToLower(input) {
	case ":q":
		os.Exit(0)
	case ":clear":
		fmt.Print("\033[H\033[2J")
		return nil
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
		return fmt.Errorf("multiple assignment")
	}

	left, err := calculon.Parse(strings.TrimSpace(vars[0]))
	if err != nil {
		return err
	}

	right, err := calculon.Parse(strings.TrimSpace(vars[1]))
	if err != nil {
		return err
	}

	switch left := left.(type) {
	case calculon.Variable:
		num, ok := right.(calculon.Number)
		if !ok {
			return fmt.Errorf("invalid right operand type: %T", right)
		}

		ctx.SetVar(left.Name, num.Value)
	case calculon.FunctionCall:
		var pnames []string
		for _, arg := range left.Args {
			vararg, ok := arg.(calculon.Variable)
			if !ok {
				return fmt.Errorf("unsupported function argument: %T", arg)
			}

			pnames = append(pnames, vararg.Name)
		}

		fnCtx := NewForkCtx(ctx)
		ctx.SetFunc(left.Name, func(args []float64) (float64, error) {
			if len(args) != len(pnames) {
				return 0, fmt.Errorf("%s(): bad params count (want %d, got %d)", left.Name, len(pnames), len(args))
			}

			for i, paramName := range pnames {
				fnCtx.SetVar(paramName, args[i])
			}

			return right.Eval(fnCtx)
		})
	default:
		return fmt.Errorf("invalid left operand type: %T", left)
	}

	return nil
}
