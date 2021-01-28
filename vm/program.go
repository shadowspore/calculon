package vm

type Program struct {
	Opcodes []Opcode
	Values  []float64
	Vars    []string
	Funcs   []FuncMeta
}

type FuncMeta struct {
	Name      string
	ArgsCount int
}
