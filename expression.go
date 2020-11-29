package calculon

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Expression interface {
	Eval(ctx EvalContext) (float64, error)
	String() string
}

type Number struct {
	Value float64
}

func (c Number) Eval(ctx EvalContext) (float64, error) {
	return c.Value, nil
}

func (c Number) String() string {
	return strconv.FormatFloat(c.Value, 'g', 10, 64)
}

type Addition struct {
	Left, Right Expression
}

func (add Addition) Eval(ctx EvalContext) (float64, error) {
	l, err := add.Left.Eval(ctx)
	if err != nil {
		return 0, err
	}

	r, err := add.Right.Eval(ctx)
	if err != nil {
		return 0, err
	}

	return l + r, nil
}

func (add Addition) String() string {
	return add.Left.String() + " + " + add.Right.String()
}

type Subtraction struct {
	Left, Right Expression
}

func (sub Subtraction) Eval(ctx EvalContext) (float64, error) {
	l, err := sub.Left.Eval(ctx)
	if err != nil {
		return 0, err
	}

	r, err := sub.Right.Eval(ctx)
	if err != nil {
		return 0, err
	}

	return l - r, nil
}

func (sub Subtraction) String() string {
	return sub.Left.String() + " - " + sub.Right.String()
}

type Multiplication struct {
	Left, Right Expression
}

func (mul Multiplication) Eval(ctx EvalContext) (float64, error) {
	l, err := mul.Left.Eval(ctx)
	if err != nil {
		return 0, err
	}

	r, err := mul.Right.Eval(ctx)
	if err != nil {
		return 0, err
	}

	return l * r, nil
}

func (mul Multiplication) String() string {
	return mul.Left.String() + " * " + mul.Right.String()
}

type Division struct {
	Left, Right Expression
}

func (div Division) Eval(ctx EvalContext) (float64, error) {
	l, err := div.Left.Eval(ctx)
	if err != nil {
		return 0, err
	}

	r, err := div.Right.Eval(ctx)
	if err != nil {
		return 0, err
	}

	if r == 0 {
		return 0, fmt.Errorf("divide by zero")
	}

	return l / r, nil
}

func (div Division) String() string {
	return div.Left.String() + " / " + div.Right.String()
}

type Modulo struct {
	Left, Right Expression
}

func (mod Modulo) Eval(ctx EvalContext) (float64, error) {
	l, err := mod.Left.Eval(ctx)
	if err != nil {
		return 0, err
	}

	r, err := mod.Right.Eval(ctx)
	if err != nil {
		return 0, err
	}

	return math.Mod(l, r), nil
}

func (mod Modulo) String() string {
	return mod.Left.String() + " % " + mod.Right.String()
}

type Negation struct {
	Expr Expression
}

func (neg Negation) Eval(ctx EvalContext) (float64, error) {
	num, err := neg.Expr.Eval(ctx)
	if err != nil {
		return 0, err
	}

	return -num, err
}

func (neg Negation) String() string {
	return "-" + neg.Expr.String()
}

type Exponentiation struct {
	Num, Power Expression
}

func (exp Exponentiation) Eval(ctx EvalContext) (float64, error) {
	n, err := exp.Num.Eval(ctx)
	if err != nil {
		return 0, err
	}

	pow, err := exp.Power.Eval(ctx)
	if err != nil {
		return 0, err
	}

	return math.Pow(n, pow), nil
}

func (exp Exponentiation) String() string {
	return exp.Num.String() + "^" + exp.Power.String()
}

type Variable struct {
	Name string
}

func (vb Variable) Eval(ctx EvalContext) (float64, error) {
	value, found := ctx.LookupVar(vb.Name)
	if !found {
		return 0, fmt.Errorf("variable not specified: %s", vb)
	}

	return value, nil
}

func (vb Variable) String() string {
	return vb.Name
}

type FunctionCall struct {
	Name string
	Args []Expression
}

func (call FunctionCall) Eval(ctx EvalContext) (float64, error) {
	fn, found := ctx.LookupFunc(call.Name)
	if !found {
		return 0, fmt.Errorf("function not specified: %s", call.Name)
	}

	args := make([]float64, 0, len(call.Args))
	for _, arg := range call.Args {
		n, err := arg.Eval(ctx)
		if err != nil {
			return 0, err
		}

		args = append(args, n)
	}

	return fn(args)
}

func (call FunctionCall) String() string {
	var args []string
	for _, arg := range call.Args {
		args = append(args, arg.String())
	}

	return call.Name + "(" + strings.Join(args, ", ") + ")"
}

type Parentheses struct {
	Expr Expression
}

func (par Parentheses) Eval(ctx EvalContext) (float64, error) {
	return par.Expr.Eval(ctx)
}

func (par Parentheses) String() string {
	return "(" + par.Expr.String() + ")"
}

// type ConditionalOp struct {
// 	Op          string
// 	Left, Right Expression
// }

// type UnaryInfixOp struct {
// 	Op   string
// 	Expr Expression
// }

// type UnaryPostfixOp struct {
// 	Op   string
// 	Expr Expression
// }

// type BinaryOp struct {
// 	Op          string
// 	Left, Right Expression
// }

// type TernaryOp struct {
// 	Op              string
// 	Cond            ConditionalOp
// 	OnTrue, OnFalse Expression
// }
