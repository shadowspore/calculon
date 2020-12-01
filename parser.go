package calculon

import (
	"fmt"
	"strconv"

	"github.com/zweihander/calculon/lexer"
)

type recursiveDescent struct {
	lexer *lexer.Lexer
}

func newRecursiveDescent(input string) *recursiveDescent {
	return &recursiveDescent{
		lexer: lexer.New(input),
	}
}

func (p *recursiveDescent) parse() (Expression, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	if next := p.lexer.Next(); next.Kind != lexer.EOF {
		return nil, fmt.Errorf("unexpected token: %s", next)
	}

	return expr, nil
}

func (p *recursiveDescent) parseExpr() (Expression, error) {
	expr, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for {
		next := p.lexer.Ahead().Kind
		if next == lexer.Plus || next == lexer.Minus {
			_ = p.lexer.Next()
			right, err := p.parseTerm()
			if err != nil {
				return nil, err
			}

			expr = BinaryOp{
				Op:    next.String(),
				Left:  expr,
				Right: right,
			}

			continue
		}

		return expr, nil
	}
}

func (p *recursiveDescent) parseTerm() (Expression, error) {
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for {
		next := p.lexer.Ahead().Kind
		if next == lexer.Asterisk || next == lexer.Slash || next == lexer.Percent {
			_ = p.lexer.Next()
			right, err := p.parseFactor()
			if err != nil {
				return nil, err
			}

			left = BinaryOp{
				Op:    next.String(),
				Left:  left,
				Right: right,
			}

			continue
		}

		return left, nil
	}
}

func (p *recursiveDescent) parseFactor() (expr Expression, err error) {
	tok := p.lexer.Next()
	switch tok.Kind {
	case lexer.Plus:
		return p.parseFactor()
	case lexer.Minus:
		expr, err := p.parseFactor()
		if err != nil {
			return nil, err
		}

		return UnaryOp{Op: "-", Expr: expr}, nil
	case lexer.OpenParen:
		expr, err = p.parseExpr()
		if err != nil {
			return nil, err
		}

		if next := p.lexer.Next(); next.Kind != lexer.CloseParen {
			return nil, fmt.Errorf("expected: ')'")
		}
	case lexer.Number:
		num, err := strconv.ParseFloat(tok.Value, 64)
		if err != nil {
			return nil, err
		}

		expr = Number{Value: num}
	case lexer.Ident:
		// it's a function?
		if p.lexer.Eat(lexer.OpenParen) {
			args, err := p.parseArgs()
			if err != nil {
				return nil, err
			}

			expr = FunctionCall{
				Name: tok.Value,
				Args: args,
			}

			break
		}

		expr = Variable{Name: tok.Value}
	default:
		return nil, fmt.Errorf("unexpected token: %s", tok)
	}

	if p.lexer.Eat(lexer.Caret) {
		power, err := p.parseFactor()
		if err != nil {
			return nil, err
		}

		return BinaryOp{
			Op:    "^",
			Left:  expr,
			Right: power,
		}, nil
	}

	return
}

func (p *recursiveDescent) parseArgs() ([]Expression, error) {
	var args []Expression
	for !p.lexer.Eat(lexer.CloseParen) {
		if len(args) > 0 {
			if !p.lexer.Eat(lexer.Comma) {
				return nil, fmt.Errorf("expected ','")
			}
		}

		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		args = append(args, expr)
	}

	return args, nil
}

func Parse(input string) (Expression, error) {
	return newRecursiveDescent(input).parse()
}
