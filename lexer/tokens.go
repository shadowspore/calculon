package lexer

type Token struct {
	Kind  Kind
	Value string
}

func (tok Token) String() string {
	kind := "Kind:" + tok.Kind.String()
	if tok.Value != "" {
		return kind + "(Value:" + tok.Value + ")"
	}

	return kind
}

type Kind byte

func (k Kind) String() string {
	switch k {
	case EOF:
		return "EOF"
	case Unexpected:
		return "Unexpected"
	case Plus:
		return "+"
	case Minus:
		return "-"
	case Asterisk:
		return "*"
	case Slash:
		return "/"
	case Percent:
		return "%"
	case Caret:
		return "^"
	case OpenParen:
		return "("
	case CloseParen:
		return ")"
	case Comma:
		return ","
	case Ident:
		return "Ident"
	case Number:
		return "Number"
	default:
		return "Unknown kind: " + string(k)
	}
}

const (
	EOF        Kind = iota
	Unexpected      //
	Plus            // +
	Minus           // -
	Asterisk        // *
	Slash           // /
	Percent         // %
	Caret           // ^
	OpenParen       // (
	CloseParen      // )
	Comma           // ,
	Ident           // foo
	Number          // 123
)
