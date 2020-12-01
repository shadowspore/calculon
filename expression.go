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

type BinaryOp struct {
	Op    string
	Left  Expression
	Right Expression
}

func (binary BinaryOp) Eval(ctx EvalContext) (float64, error) {
	l, err := binary.Left.Eval(ctx)
	if err != nil {
		return 0, err
	}

	r, err := binary.Right.Eval(ctx)
	if err != nil {
		return 0, err
	}

	switch binary.Op {
	case "+":
		return l + r, nil
	case "-":
		return l - r, nil
	case "*":
		return l * r, nil
	case "/":
		if r == 0 {
			return 0, fmt.Errorf("divide by zero")
		}

		return l / r, nil
	case "%":
		return math.Mod(l, r), nil
	case "^":
		return math.Pow(l, r), nil
	default:
		return 0, fmt.Errorf("unexpected binary op: %s", binary.Op)
	}
}

func (binary BinaryOp) String() string {
	currOp := precedense[binary.Op]

	left := binary.Left.String()
	if lbin, ok := binary.Left.(BinaryOp); ok {
		leftOp := precedense[lbin.Op]
		if currOp.prec > leftOp.prec || (currOp.prec == leftOp.prec && currOp.rightAssoc) {
			left = "(" + left + ")"
		}
	}

	right := binary.Right.String()
	if rbin, ok := binary.Right.(BinaryOp); ok {
		rightOp := precedense[rbin.Op]
		if currOp.prec > rightOp.prec ||
			(currOp.prec == rightOp.prec && !rightOp.rightAssoc && !currOp.rightAssoc) {
			right = "(" + right + ")"
		}
	}

	if binary.Op == "^" {
		return left + binary.Op + right
	}

	return left + " " + binary.Op + " " + right
}

type UnaryOp struct {
	Op        string
	Expr      Expression
	IsPostfix bool
}

func (unary UnaryOp) Eval(ctx EvalContext) (float64, error) {
	val, err := unary.Expr.Eval(ctx)
	if err != nil {
		return 0, err
	}

	switch unary.Op {
	case "-":
		return -val, nil
	default:
		return 0, fmt.Errorf("unexpected unary op: %s", unary.Op)
	}
}

func (unary UnaryOp) String() string {
	expString := unary.Expr.String()
	if _, ok := unary.Expr.(BinaryOp); ok {
		expString = "(" + expString + ")"
	}

	if unary.IsPostfix {
		return expString + unary.Op
	}

	return unary.Op + expString
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
