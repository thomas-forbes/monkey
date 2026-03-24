package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"monkey/object"
	"monkey/runner"
	"os"
)

const prompt = ">> "

func main() {
	engineFlag := flag.String("engine", "vm", "use 'vm' or 'eval'")
	flag.Parse()

	engine, err := runner.ParseEngine(*engineFlag)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	if flag.Arg(0) == "" {
		startRepl(engine, os.Stdin, os.Stdout)
	} else if len(os.Args) > 1 {
		startFile(engine, flag.Arg(0), os.Stdout)
	} else {
		fmt.Println("Usage: monkey file.monkey")
	}
}

func startRepl(engine runner.Engine, in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	session := runner.NewSession(engine)
	fmt.Fprintf(out, "Running Monkey REPL with %s engine\n", engine)

	for {
		fmt.Fprintf(out, prompt)
		if !scanner.Scan() {
			return
		}

		result, _ := runner.RunProgram(engine, scanner.Text(), session)

		if result != nil {
			io.WriteString(out, result.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func startFile(engine runner.Engine, fileName string, out io.Writer) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(out, "Couldn't read file: %s\n", err)
		return
	}

	session := runner.NewSession(engine)
	result, duration := runner.RunProgram(engine, string(data), session)
	if result == nil {
		result = object.NULL
	}

	fmt.Fprintf(out, "engine=%s, result=%s, duration=%s\n", engine, result.Inspect(), duration)
}
