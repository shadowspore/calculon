package vm

type Opcode byte

const (
	_ Opcode = iota
	Num
	Add
	Sub
	Mult
	Div
	Mod
	Exp
	Neg
	Var
	Func
)

func (o Opcode) String() string {
	switch o {
	case Num:
		return "Num"
	case Add:
		return "Add"
	case Sub:
		return "Sub"
	case Mult:
		return "Mult"
	case Div:
		return "Div"
	case Mod:
		return "Mod"
	case Exp:
		return "Exp"
	case Neg:
		return "Neg"
	case Var:
		return "Var"
	case Func:
		return "Func"
	default:
		return "UnknownOpcode"
	}
}
