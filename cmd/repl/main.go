package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/xjem/calculon"
	"github.com/xjem/calculon/internal/repl"
)

func main() {
	var (
		reader = bufio.NewReader(os.Stdin)
		repl   = repl.New(calculon.MathContext())
	)

	for {
		fmt.Print(">> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		if err := func(input string) error {
			switch input {
			case ":q":
				os.Exit(0)
			case ":clear":
				fmt.Print("\033[H\033[2J")
				return nil
			}

			if strings.Contains(input, "=") {
				return repl.Define(input)
			}

			result, err := repl.Eval(input)
			if err != nil {
				return err
			}

			fmt.Println(result)
			return nil
		}(strings.TrimSpace(input)); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
		}
	}
}
