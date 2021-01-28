package ast

import (
	"strconv"
	"strings"
)

type Node interface {
	String() string
}

type Number struct {
	Value float64
}

func (c Number) String() string {
	return strconv.FormatFloat(c.Value, 'g', 10, 64)
}

type BinaryOp struct {
	Op    string
	Left  Node
	Right Node
}

func (binary BinaryOp) String() string {
	if binary.Op == "^" {
		return binary.Left.String() + binary.Op + binary.Right.String()
	}

	return binary.Left.String() + " " + binary.Op + " " + binary.Right.String()
}

type UnaryOp struct {
	Op        string
	Operand   Node
	IsPostfix bool
}

func (unary UnaryOp) String() string {
	if unary.IsPostfix {
		return unary.Operand.String() + unary.Op
	}

	return unary.Op + unary.Operand.String()
}

type Parentheses struct {
	Inner Node
}

func (paren Parentheses) String() string {
	return "(" + paren.Inner.String() + ")"
}

type Variable struct {
	Name string
}

func (vb Variable) String() string {
	return vb.Name
}

type FunctionCall struct {
	Name string
	Args []Node
}

func (call FunctionCall) String() string {
	var args []string
	for _, arg := range call.Args {
		args = append(args, arg.String())
	}

	return call.Name + "(" + strings.Join(args, ", ") + ")"
}
