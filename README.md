# calculon

Simple math expression library

## Example

```go
package main

import (
	"fmt"

	"github.com/zweihander/calculon"
)

func main() {
	expr, err := calculon.Parse("sin(Pi/x)")
	if err != nil {
		panic(err)
	}

	// MathContext have math std variables and functions
	// such as sin, cos, Pi, etc
	ctx := calculon.MathContext()

	ctx.SetVar("x", 2) // set variable
	result, err := expr.Eval(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(result) // 1
}

```

## REPL

Interactive calculator in ```cmd/repl```

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

## About implementation

### Parsing
Parser uses [recursive descent algorithm](https://en.wikipedia.org/wiki/Recursive_descent_parser) to parse [expression grammar](GRAMMAR.md).\
I didn't use shunting-yard algorithm because it uses [some hacky solutions](https://stackoverflow.com/a/17132657) to handle unary minus operator, which may conflict with user-defined variables or functions.

### Interpreter
It's a tree-walking interpreter.
Each AST node have ```Eval()``` function.

## Articles
1. [Java expression evaluator](https://stackoverflow.com/a/26227947)
2. [Writing a Simple Math Expression Engine in C#](https://medium.com/@toptensoftware/writing-a-simple-math-expression-engine-in-c-d414de18d4ce)
3. [Pretty Printing AST with Minimal Parentheses](https://stackoverflow.com/questions/13708837/pretty-printing-ast-with-minimal-parentheses)
4. [Shunting-yard algorithm](https://rosettacode.org/wiki/Parsing/Shunting-yard_algorithm#Go)