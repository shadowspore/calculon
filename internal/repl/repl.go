package repl

import (
	"fmt"
	"strings"

	"github.com/xjem/calculon"
)

type Repl struct {
	globalScope *Scope
}

func New(std calculon.EvalContext) *Repl {
	return &Repl{
		globalScope: NewScope(std),
	}
}

func (r *Repl) Eval(input string) (float64, error) {
	expr, err := calculon.Parse(input)
	if err != nil {
		return 0, fmt.Errorf("parse: %w", err)
	}

	return expr.Eval(r.globalScope)
}

func (r *Repl) Define(input string) error {
	vars := strings.Split(input, "=")
	if len(vars) != 2 {
		return fmt.Errorf("multiple assignment")
	}

	definition, err := calculon.Parse(strings.TrimSpace(vars[0]))
	if err != nil {
		return err
	}

	body, err := calculon.Parse(strings.TrimSpace(vars[1]))
	if err != nil {
		return err
	}

	switch definition := definition.(type) {
	case calculon.Variable:
		num, ok := body.(calculon.Number)
		if !ok {
			return fmt.Errorf("invalid variable type: %T", body)
		}

		r.globalScope.SetVar(definition.Name, num.Value)

	case calculon.FunctionCall:
		var requiredArgs []string
		for _, arg := range definition.Args {
			vararg, ok := arg.(calculon.Variable)
			if !ok {
				return fmt.Errorf("unsupported function argument: %T", arg)
			}

			requiredArgs = append(requiredArgs, vararg.Name)
		}

		fnScope := NewScope(r.globalScope)
		r.globalScope.SetFunc(definition.Name, func(args []float64) (float64, error) {
			if len(args) != len(requiredArgs) {
				return 0, fmt.Errorf("%s: bad params count (want %d, got %d)", definition.String(), len(requiredArgs), len(args))
			}

			for i, paramName := range requiredArgs {
				fnScope.vars[paramName] = args[i]
			}

			return body.Eval(fnScope)
		})
	default:
		return fmt.Errorf("invalid definition type: %T", definition)
	}

	return nil
}
