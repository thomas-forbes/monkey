package main

import (
	"fmt"
	"monkey/runner"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		runner.StartRepl(os.Stdin, os.Stdout)
	} else if len(os.Args) == 2 {
		runner.StartFile(os.Args[1])
	} else {
		fmt.Println("Usage: monkey file.monkey")
	}
}
