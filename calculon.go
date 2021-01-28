package calculon

import (
	"github.com/zweihander/calculon/compiler"
	"github.com/zweihander/calculon/parser"
	"github.com/zweihander/calculon/vm"
)

func NewVM() *vm.VM {
	return vm.New(vm.Config{
		StackPoolSize: 10,
		StackCapacity: 10,
	})
}

func Compile(input string) (vm.Program, error) {
	node, err := parser.New(input).Parse()
	if err != nil {
		return vm.Program{}, err
	}

	return compiler.Compile(node)
}
