package compiler

import (
	"fmt"

	"github.com/zweihander/calculon/ast"
	"github.com/zweihander/calculon/vm"
)

type compiler struct {
	ops   []vm.Opcode
	vals  []float64
	vars  []string
	funcs []vm.FuncMeta
}

func Compile(root ast.Node) (vm.Program, error) {
	var c compiler
	if err := c.compile(root); err != nil {
		return vm.Program{}, err
	}

	return vm.Program{
		Opcodes: c.ops,
		Values:  c.vals,
		Vars:    c.vars,
		Funcs:   c.funcs,
	}, nil
}

func (c *compiler) compile(node ast.Node) error {
	switch n := node.(type) {
	case ast.BinaryOp:
		return c.binaryOp(n)
	case ast.UnaryOp:
		return c.unaryOp(n)
	case ast.Number:
		c.vals = append(c.vals, n.Value)
		c.ops = append(c.ops, vm.Num)
		return nil
	case ast.Parentheses:
		return c.compile(n.Inner)
	case ast.Variable:
		c.vars = append(c.vars, n.Name)
		c.ops = append(c.ops, vm.Var)
		return nil
	case ast.FunctionCall:
		return c.funcCall(n)
	default:
		return fmt.Errorf("unexpected ast node: %T", n)
	}
}

func (c *compiler) binaryOp(op ast.BinaryOp) error {
	if err := c.compile(op.Left); err != nil {
		return err
	}

	if err := c.compile(op.Right); err != nil {
		return err
	}

	switch op.Op {
	case "+":
		c.ops = append(c.ops, vm.Add)
	case "-":
		c.ops = append(c.ops, vm.Sub)
	case "*":
		c.ops = append(c.ops, vm.Mult)
	case "/":
		c.ops = append(c.ops, vm.Div)
	case "%":
		c.ops = append(c.ops, vm.Mod)
	case "^":
		c.ops = append(c.ops, vm.Exp)
	default:
		return fmt.Errorf("unknown binary op: %s", op.Op)
	}

	return nil
}

func (c *compiler) unaryOp(op ast.UnaryOp) error {
	if err := c.compile(op.Operand); err != nil {
		return err
	}

	switch op.Op {
	case "-":
		c.ops = append(c.ops, vm.Neg)
	default:
		return fmt.Errorf("unknown unary op: %s", op.Op)
	}

	return nil
}

func (c *compiler) funcCall(fn ast.FunctionCall) error {
	c.funcs = append(c.funcs, vm.FuncMeta{
		Name:      fn.Name,
		ArgsCount: len(fn.Args),
	})

	for _, arg := range fn.Args {
		if err := c.compile(arg); err != nil {
			return err
		}
	}

	c.ops = append(c.ops, vm.Func)

	return nil
}
