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

func (p *recursiveDescent) parseExpr() (Expression, error) {
	expr, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for {
		if next := p.lexer.Ahead().Kind; next == lexer.CloseParen || next == lexer.Comma {
			return expr, nil
		}

		tok := p.lexer.Next()
		switch tok.Kind {
		case lexer.Plus:
			plusExpr, err := p.parseTerm()
			if err != nil {
				return nil, err
			}

			expr = Addition{Left: expr, Right: plusExpr}
		case lexer.Minus:
			minusExpr, err := p.parseTerm()
			if err != nil {
				return nil, err
			}

			expr = Subtraction{Left: expr, Right: minusExpr}
		case lexer.EOF:
			return expr, nil
		default:
			return nil, fmt.Errorf("unexpected token: %s, pos: %d", tok, p.lexer.Pos())
		}
	}
}

func (p *recursiveDescent) parseTerm() (Expression, error) {
	expr, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for {
		if p.lexer.Eat(lexer.Asterisk) {
			mulExpr, err := p.parseFactor()
			if err != nil {
				return nil, err
			}

			expr = Multiplication{Left: expr, Right: mulExpr}
			continue
		}

		if p.lexer.Eat(lexer.Slash) {
			divExpr, err := p.parseFactor()
			if err != nil {
				return nil, err
			}

			expr = Division{Left: expr, Right: divExpr}
			continue
		}

		if p.lexer.Eat(lexer.Percent) {
			prcExpr, err := p.parseFactor()
			if err != nil {
				return nil, err
			}

			expr = Modulo{Left: expr, Right: prcExpr}
			continue
		}

		return expr, nil
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

		return Negation{Expr: expr}, nil
	case lexer.OpenParen:
		expr, err = p.parseExpr()
		if err != nil {
			return nil, err
		}

		if next := p.lexer.Next(); next.Kind != lexer.CloseParen {
			return nil, fmt.Errorf("expected: ')'")
		}

		expr = Parentheses{Expr: expr}
	case lexer.Number:
		num, err := strconv.ParseFloat(tok.Value, 64)
		if err != nil {
			return nil, err
		}

		expr = Number{Value: num}
	case lexer.Ident:
		// it's a function?
		if !p.lexer.Eat(lexer.OpenParen) { // no
			expr = Variable{Name: tok.Value}
			break
		}

		// it's a function without args?
		if p.lexer.Eat(lexer.CloseParen) { // yes
			expr = FunctionCall{
				Name: tok.Value,
				Args: nil,
			}

			break
		}

		var args []Expression

		for {
			arg, err := p.parseExpr()
			if err != nil {
				return nil, err
			}

			args = append(args, arg)
			if p.lexer.Eat(lexer.Comma) {
				continue
			}

			break
		}

		if p.lexer.Next().Kind != lexer.CloseParen {
			return nil, fmt.Errorf("expected ')'")
		}

		expr = FunctionCall{
			Name: tok.Value,
			Args: args,
		}

	default:
		return nil, fmt.Errorf("unexpected token: %s", tok)
	}

	if p.lexer.Eat(lexer.Caret) {
		power, err := p.parseFactor()
		if err != nil {
			return nil, err
		}

		return Exponentiation{
			Num:   expr,
			Power: power,
		}, nil
	}

	return
}

func Parse(input string) (Expression, error) {
	return newRecursiveDescent(input).parseExpr()
}
