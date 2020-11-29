package lexer

import (
	"unicode"
)

type Lexer struct {
	source []rune
	pos    int
	end    bool
}

func (l *Lexer) Pos() int {
	return l.pos
}

func New(input string) *Lexer {
	return &Lexer{
		source: []rune(input),
		end:    len(input) == 0,
	}
}

func (l *Lexer) current() rune {
	return l.source[l.pos]
}

func (l *Lexer) next() bool {
	if l.pos == len(l.source)-1 {
		l.end = true
		return false
	}

	l.pos++
	return true
}

func (l *Lexer) Next() Token {
	if l.end {
		return Token{Kind: EOF}
	}

	for unicode.IsSpace(l.current()) {
		if !l.next() {
			return Token{Kind: EOF}
		}
	}

	switch l.current() {
	case '+':
		l.next()
		return Token{Kind: Plus}
	case '-':
		l.next()
		return Token{Kind: Minus}
	case '*':
		l.next()
		return Token{Kind: Asterisk}
	case '/':
		l.next()
		return Token{Kind: Slash}
	case '%':
		l.next()
		return Token{Kind: Percent}
	case '^':
		l.next()
		return Token{Kind: Caret}
	case '(':
		l.next()
		return Token{Kind: OpenParen}
	case ')':
		l.next()
		return Token{Kind: CloseParen}
	case ',':
		l.next()
		return Token{Kind: Comma}
	}

	if unicode.IsDigit(l.current()) {
		var num []rune

		haveDecimalPoint := false
		for unicode.IsDigit(l.current()) || (!haveDecimalPoint && l.current() == '.') {
			num = append(num, l.current())
			haveDecimalPoint = l.current() == '.'

			if !l.next() { // EOF
				break
			}
		}

		return Token{
			Kind:  Number,
			Value: string(num),
		}
	}

	if unicode.IsLetter(l.current()) || l.current() == '_' {
		var ident []rune

		for unicode.IsLetter(l.current()) || l.current() == '_' {
			ident = append(ident, l.current())

			if !l.next() { // EOF
				break
			}
		}

		return Token{
			Kind:  Ident,
			Value: string(ident),
		}
	}

	return Token{
		Kind:  Unexpected,
		Value: string(l.current()),
	}
}

func (l *Lexer) Eat(expect Kind) bool {
	prevPos := l.pos
	prevEnd := l.end

	tok := l.Next()
	if tok.Kind != expect {
		l.pos = prevPos
		l.end = prevEnd
		return false
	}

	return true
}

func (l *Lexer) Ahead() Token {
	prevPos := l.pos
	prevEnd := l.end

	tok := l.Next()
	l.pos = prevPos
	l.end = prevEnd

	return tok
}
