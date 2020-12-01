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

func (p *recursiveDescent) parseFactor() (Expression, error) {
	if p.lexer.Eat(lexer.Plus) {
		return p.parseFactor()
	}

	if p.lexer.Eat(lexer.Minus) {
		expr, err := p.parseFactor()
		if err != nil {
			return nil, err
		}

		return UnaryOp{Op: "-", Expr: expr}, nil
	}

	expr, err := p.parsePrimary()
	if err != nil {
		return nil, err
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

	return expr, nil
}

func (p *recursiveDescent) parsePrimary() (Expression, error) {
	tok := p.lexer.Next()
	switch tok.Kind {
	case lexer.OpenParen:
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		if !p.lexer.Eat(lexer.CloseParen) {
			return nil, fmt.Errorf("expected: ')'")
		}

		return expr, nil
	case lexer.Number:
		num, err := strconv.ParseFloat(tok.Value, 64)
		if err != nil {
			return nil, err
		}

		return Number{Value: num}, nil
	case lexer.Ident:
		// it's a function?
		if p.lexer.Eat(lexer.OpenParen) {
			args, err := p.parseArgs()
			if err != nil {
				return nil, err
			}

			return FunctionCall{
				Name: tok.Value,
				Args: args,
			}, nil
		}

		return Variable{Name: tok.Value}, nil
	default:
		return nil, fmt.Errorf("unexpected token: %s", tok)
	}
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
