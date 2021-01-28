package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/zweihander/calculon"
	"github.com/zweihander/calculon/ast"
	"github.com/zweihander/calculon/compiler"
	"github.com/zweihander/calculon/parser"
	"github.com/zweihander/calculon/vm"
)

type ScopedEnv struct {
	parent *calculon.Env
	*calculon.Env
}

func (e *ScopedEnv) LookupVar(name string) (float64, bool) {
	val, found := e.Env.LookupVar(name)
	if found {
		return val, true
	}

	return e.parent.LookupVar(name)
}

func (e *ScopedEnv) LookupFunc(name string) (calculon.Function, bool) {
	fn, found := e.Env.LookupFunc(name)
	if found {
		return fn, true
	}

	return e.parent.LookupFunc(name)
}

func NewScopedEnv(parent *calculon.Env) *ScopedEnv {
	return &ScopedEnv{
		parent: parent,
		Env:    calculon.NewEnv(),
	}
}

func main() {
	r := bufio.NewReader(os.Stdin)
	repl := &Repl{
		global: calculon.MathEnv(),
		vm:     vm.New(vm.Config{}),
	}

	for {
		fmt.Print(">> ")
		input, err := r.ReadString('\n')
		if err != nil {
			panic(err)
		}

		if err := repl.Exec(strings.TrimSpace(input)); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			continue
		}
	}
}

type Repl struct {
	global *calculon.Env
	vm     *vm.VM
}

func (repl *Repl) Exec(input string) error {
	if strings.Contains(input, "=") {
		if err := repl.define(input); err != nil {
			return fmt.Errorf("assign: %w", err)
		}

		return nil
	}

	if strings.ToLower(input) == ":q" {
		os.Exit(0)
	}

	program, err := calculon.Compile(input)
	if err != nil {
		return fmt.Errorf("compile: %w", err)
	}

	result, err := repl.vm.Run(program, repl.global)
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}

	fmt.Println(result)
	return nil
}

// hacky hack
func (repl *Repl) define(input string) error {
	vars := strings.Split(input, "=")
	if len(vars) != 2 {
		return fmt.Errorf("multiple assignment")
	}

	left, err := parser.Parse(strings.TrimSpace(vars[0]))
	if err != nil {
		return err
	}

	right, err := parser.Parse(strings.TrimSpace(vars[1]))
	if err != nil {
		return err
	}

	switch left := left.(type) {
	case ast.Variable:
		num, ok := right.(ast.Number)
		if !ok {
			return fmt.Errorf("invalid right operand type: %T", right)
		}

		repl.global.SetVar(left.Name, num.Value)
	case ast.FunctionCall:
		var pnames []string
		for _, arg := range left.Args {
			vararg, ok := arg.(ast.Variable)
			if !ok {
				return fmt.Errorf("unsupported function argument: %T", arg)
			}

			pnames = append(pnames, vararg.Name)
		}

		fnProgram, err := compiler.Compile(right)
		if err != nil {
			return fmt.Errorf("compile: %w", err)
		}

		fnEnv := NewScopedEnv(repl.global)
		fnEnv.SetFunc(left.Name, func(args []float64) (float64, error) {
			if len(args) != len(pnames) {
				return 0, fmt.Errorf("%s(): bad params count (want %d, got %d)", left.Name, len(pnames), len(args))
			}

			for i, paramName := range pnames {
				fnEnv.SetVar(paramName, args[i])
			}

			return repl.vm.Run(fnProgram, fnEnv)
		})
	default:
		return fmt.Errorf("invalid left operand type: %T", left)
	}

	return nil
}
