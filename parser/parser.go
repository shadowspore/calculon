package parser

import (
	"fmt"
	"strconv"

	"github.com/zweihander/calculon/ast"
	"github.com/zweihander/calculon/parser/lexer"
)

type Parser struct {
	lexer *lexer.Lexer
}

func New(input string) *Parser {
	return &Parser{
		lexer: lexer.New(input),
	}
}

func (p *Parser) Parse() (ast.Node, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	if next := p.lexer.Next(); next.Kind != lexer.EOF {
		return nil, fmt.Errorf("unexpected token: %s", next)
	}

	return expr, nil
}

func (p *Parser) parseExpr() (ast.Node, error) {
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

			expr = ast.BinaryOp{
				Op:    next.String(),
				Left:  expr,
				Right: right,
			}

			continue
		}

		return expr, nil
	}
}

func (p *Parser) parseTerm() (ast.Node, error) {
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

			left = ast.BinaryOp{
				Op:    next.String(),
				Left:  left,
				Right: right,
			}

			continue
		}

		return left, nil
	}
}

func (p *Parser) parseFactor() (ast.Node, error) {
	if p.lexer.Eat(lexer.Plus) {
		return p.parseFactor()
	}

	if p.lexer.Eat(lexer.Minus) {
		expr, err := p.parseFactor()
		if err != nil {
			return nil, err
		}

		return ast.UnaryOp{Op: "-", Operand: expr}, nil
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

		return ast.BinaryOp{
			Op:    "^",
			Left:  expr,
			Right: power,
		}, nil
	}

	return expr, nil
}

func (p *Parser) parsePrimary() (ast.Node, error) {
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

		return ast.Parentheses{Inner: expr}, nil
	case lexer.Number:
		num, err := strconv.ParseFloat(tok.Value, 64)
		if err != nil {
			return nil, err
		}

		return ast.Number{Value: num}, nil
	case lexer.Ident:
		// it's a function?
		if p.lexer.Eat(lexer.OpenParen) {
			args, err := p.parseArgs()
			if err != nil {
				return nil, err
			}

			return ast.FunctionCall{
				Name: tok.Value,
				Args: args,
			}, nil
		}

		return ast.Variable{Name: tok.Value}, nil
	default:
		return nil, fmt.Errorf("unexpected token: %s", tok)
	}
}

func (p *Parser) parseArgs() ([]ast.Node, error) {
	var args []ast.Node
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

func Parse(input string) (ast.Node, error) {
	return New(input).Parse()
}
