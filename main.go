package main

import (
	"flag"
	"fmt"
	"monkey/runner"
	"os"
)

func main() {
	flag.Parse()

	engine, err := runner.GetEngine()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	if len(os.Args) == 1 {
		runner.StartRepl(engine, os.Stdin, os.Stdout)
	} else if len(os.Args) > 1 {
		runner.StartFile(engine, os.Args[len(os.Args)-1])
	} else {
		fmt.Println("Usage: monkey file.monkey")
	}
}
