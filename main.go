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
	session, err := runner.NewSession(engine)
	if err != nil {
		fmt.Fprintf(out, "Error initializing session: %s\n", err)
		return
	}
	fmt.Fprintf(out, "Running Monkey REPL with %s engine\n", engine)

	for {
		fmt.Fprintf(out, prompt)
		if !scanner.Scan() {
			return
		}

		program, err := runner.ParseCode(scanner.Text())
		if err != nil {
			fmt.Fprintf(out, "Error parsing code: %s\n", err)
			continue
		}
		result, _ := session.ExecProgram(program)

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

	program, errResult := runner.ParseCode(string(data))
	if errResult != nil {
		fmt.Fprintf(out, "Error parsing code: %s\n", errResult)
		return
	}

	session, errResult := runner.NewSession(engine)
	if errResult != nil {
		fmt.Fprintf(out, "Error initializing session: %s\n", errResult)
		return
	}
	result, duration := session.ExecProgram(program)
	if result == nil {
		result = object.NULL
	}

	fmt.Fprintf(out, "engine=%s, result=%s, duration=%s\n", engine, result.Inspect(), duration)
}
