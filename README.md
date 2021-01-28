# calculon

Simple math expression library.

Inspired by [antonmedv/expr](https://github.com/antonmedv/expr).

## Example

```go
package main

import (
	"fmt"

	"github.com/zweihander/calculon"
)

func main() {
	program, err := calculon.Compile("sin(Pi/x)")
	if err != nil {
		panic(err)
	}

	// MathEnv have math std variables and functions
	// such as sin, cos, Pi, etc
	env := calculon.MathEnv()
	env.SetVar("x", 2) // set variable

	vm := calculon.NewVM()
	result, err := vm.Run(program, env)
	if err != nil {
		panic(err)
	}

	fmt.Println(result) // 1
}
```

## REPL (cmd/repl)

```
>> 2 * -(3 + 4)
-14
>> kek = 1337
>> kek
1337
>> f(x) = x * kek
>> f(2)
2674
>> foo(x, y) = sin(f(x))^y
>> foo(1, 2)
0.9376714427474894
>> foo(f(2), 0)
1
>> ...
```